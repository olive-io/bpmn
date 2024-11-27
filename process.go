/*
   Copyright 2023 The bpmn Authors

   This library is free software; you can redistribute it and/or
   modify it under the terms of the GNU Lesser General Public
   License as published by the Free Software Foundation; either
   version 2.1 of the License, or (at your option) any later version.

   This library is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
   Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public
   License along with this library;
*/

package bpmn

import (
	"context"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type Options struct {
	ctx                            context.Context
	idGenerator                    id.IGenerator
	locator                        data.IFlowDataLocator
	eventIngress                   event.IConsumer
	eventEgress                    event.ISource
	idGeneratorBuilder             id.IGeneratorBuilder
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder
	tracer                         tracing.ITracer
}

func NewOptions(opts ...Option) *Options {
	var options Options

	for _, opt := range opts {
		opt(&options)
	}

	if options.ctx == nil {
		options.ctx = context.Background()
	}

	if options.idGeneratorBuilder == nil {
		options.idGeneratorBuilder = id.DefaultIdGeneratorBuilder
	}

	if options.eventDefinitionInstanceBuilder == nil {
		options.eventDefinitionInstanceBuilder = event.WrappingDefinitionInstanceBuilder
	}

	if options.eventIngress == nil && options.eventEgress == nil {
		fanOut := event.NewFanOut()
		options.eventIngress = fanOut
		options.eventEgress = fanOut
	}

	if options.tracer == nil {
		options.tracer = tracing.NewTracer(options.ctx)
	}

	return &options
}

// Option allows to modify configuration of
// an instance in a flexible fashion (as it's just a modification
// function)
type Option func(*Options)

// WithTracer overrides instance's tracer
func WithTracer(tracer tracing.ITracer) Option {
	return func(opt *Options) {
		opt.tracer = tracer
	}
}

// WithContext will pass a given context to a new instance
// instead of implicitly generated one
func WithContext(ctx context.Context) Option {
	return func(opt *Options) {
		opt.ctx = ctx
	}
}

func WithIdGenerator(builder id.IGeneratorBuilder) Option {
	return func(opt *Options) {
		opt.idGeneratorBuilder = builder
	}
}

func WithEventIngress(consumer event.IConsumer) Option {
	return func(opt *Options) {
		opt.eventIngress = consumer
	}
}

func WithEventEgress(source event.ISource) Option {
	return func(opt *Options) {
		opt.eventEgress = source
	}
}

func WithLocator(locator data.IFlowDataLocator) Option {
	return func(opt *Options) {
		opt.locator = locator
	}
}

func WithVariables(variables map[string]any) Option {
	return func(opt *Options) {
		if opt.locator == nil {
			opt.locator = data.NewFlowDataLocator()
		}
		for key, value := range variables {
			opt.locator.SetVariable(key, value)
		}
	}
}

func WithDataObjects(dataObjects map[string]any) Option {
	return func(opt *Options) {
		if opt.locator == nil {
			opt.locator = data.NewFlowDataLocator()
		}
		for dataObjectId, dataObject := range dataObjects {
			locator, found := opt.locator.FindIItemAwareLocator(data.LocatorObject)
			if !found {
				locator = data.NewDataObjectContainer()
				opt.locator.PutIItemAwareLocator(data.LocatorObject, locator)
			}
			container := data.NewContainer(nil)
			container.Put(dataObject)
			locator.PutItemAwareById(dataObjectId, container)
		}
	}
}

func WithEventDefinitionInstanceBuilder(builder event.IDefinitionInstanceBuilder) Option {
	return func(opt *Options) {
		opt.eventDefinitionInstanceBuilder = builder
	}
}

type Process struct {
	*Options

	Element        *schema.Process
	Definitions    *schema.Definitions
	instances      []*Instance
	subTracerMaker func() tracing.ITracer
}

func MakeProcess(element *schema.Process, definitions *schema.Definitions, opts ...Option) Process {
	options := NewOptions(opts...)

	process := Process{
		Options:     options,
		Element:     element,
		Definitions: definitions,
		instances:   make([]*Instance, 0),
	}

	ctx := process.ctx
	process.subTracerMaker = func() tracing.ITracer {
		subTracer := tracing.NewTracer(ctx)
		tracing.NewRelay(ctx, subTracer, process.tracer, func(trace tracing.ITrace) []tracing.ITrace {
			return []tracing.ITrace{ProcessTrace{
				Process: process.Element,
				Trace:   trace,
			}}
		})
		return subTracer
	}

	return process
}

func NewProcess(element *schema.Process, definitions *schema.Definitions, opts ...Option) *Process {
	process := MakeProcess(element, definitions, opts...)
	return &process
}

func (p *Process) Process() *schema.Process { return p.Element }

func (p *Process) Tracer() tracing.ITracer { return p.tracer }

func (p *Process) Instantiate(opts ...Option) (inst *Instance, err error) {
	subTracer := p.subTracerMaker()

	opts = append([]Option{
		WithIdGenerator(p.idGeneratorBuilder),
		WithEventDefinitionInstanceBuilder(p.eventDefinitionInstanceBuilder),
		WithEventEgress(p.eventEgress),
		WithEventIngress(p.eventIngress),
		WithTracer(subTracer),
	}, opts...)

	options := NewOptions(opts...)

	inst, err = NewInstance(p.Element, p.Definitions, options)
	if err != nil {
		return
	}

	return
}
