// Copyright 2023 The bpmn Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package process

import (
	"context"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tools/id"
	"github.com/olive-io/bpmn/tracing"
)

type Process struct {
	Element                        *schema.Process
	Definitions                    *schema.Definitions
	instances                      []*instance.Instance
	EventIngress                   event.IConsumer
	EventEgress                    event.ISource
	idGeneratorBuilder             id.IGeneratorBuilder
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder
	Tracer                         tracing.ITracer
	subTracerMaker                 func() tracing.ITracer
}

type Option func(context.Context, *Process) context.Context

func WithIdGenerator(builder id.IGeneratorBuilder) Option {
	return func(ctx context.Context, process *Process) context.Context {
		process.idGeneratorBuilder = builder
		return ctx
	}
}

func WithEventIngress(consumer event.IConsumer) Option {
	return func(ctx context.Context, process *Process) context.Context {
		process.EventIngress = consumer
		return ctx
	}
}

func WithEventEgress(source event.ISource) Option {
	return func(ctx context.Context, process *Process) context.Context {
		process.EventEgress = source
		return ctx
	}
}

func WithEventDefinitionInstanceBuilder(builder event.IDefinitionInstanceBuilder) Option {
	return func(ctx context.Context, process *Process) context.Context {
		process.eventDefinitionInstanceBuilder = builder
		return ctx
	}
}

// WithTracer overrides process's tracer
func WithTracer(tracer tracing.ITracer) Option {
	return func(ctx context.Context, process *Process) context.Context {
		process.Tracer = tracer
		return ctx
	}
}

// WithContext will pass a given context to a new process
// instead of implicitly generated one
func WithContext(newCtx context.Context) Option {
	return func(ctx context.Context, process *Process) context.Context {
		return newCtx
	}
}

func Make(element *schema.Process, definitions *schema.Definitions, options ...Option) Process {
	process := Process{
		Element:     element,
		Definitions: definitions,
		instances:   make([]*instance.Instance, 0),
	}

	ctx := context.Background()

	for _, option := range options {
		ctx = option(ctx, &process)
	}

	if process.idGeneratorBuilder == nil {
		process.idGeneratorBuilder = id.DefaultIdGeneratorBuilder
	}

	if process.eventDefinitionInstanceBuilder == nil {
		process.eventDefinitionInstanceBuilder = event.WrappingDefinitionInstanceBuilder
	}

	if process.EventIngress == nil && process.EventEgress == nil {
		fanOut := event.NewFanOut()
		process.EventIngress = fanOut
		process.EventEgress = fanOut
	}

	if process.Tracer == nil {
		process.Tracer = tracing.NewTracer(ctx)
	}

	process.subTracerMaker = func() tracing.ITracer {
		subTracer := tracing.NewTracer(ctx)
		tracing.NewRelay(ctx, subTracer, process.Tracer, func(trace tracing.ITrace) []tracing.ITrace {
			return []tracing.ITrace{Trace{
				Process: process.Element,
				Trace:   trace,
			}}
		})
		return subTracer
	}

	return process
}

func New(element *schema.Process, definitions *schema.Definitions, options ...Option) *Process {
	process := Make(element, definitions, options...)
	return &process
}

func (process *Process) Instantiate(options ...instance.Option) (inst *instance.Instance, err error) {
	subTracer := process.subTracerMaker()

	options = append([]instance.Option{
		instance.WithIdGenerator(process.idGeneratorBuilder),
		instance.WithEventDefinitionInstanceBuilder(process.eventDefinitionInstanceBuilder),
		instance.WithEventEgress(process.EventEgress),
		instance.WithEventIngress(process.EventIngress),
		instance.WithTracer(subTracer),
	}, options...)
	inst, err = instance.NewInstance(process.Element, process.Definitions, options...)
	if err != nil {
		return
	}

	return
}
