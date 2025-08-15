/*
Copyright 2023 The bpmn Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bpmn

import (
	"context"
	"fmt"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type EngineOptions struct {
	ctx context.Context

	idGeneratorBuilder             id.IGeneratorBuilder
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder
}

func NewEngineOptions(opts ...EngineOption) *EngineOptions {
	options := &EngineOptions{ctx: context.Background()}
	for _, opt := range opts {
		opt(options)
	}

	if options.idGeneratorBuilder == nil {
		options.idGeneratorBuilder = id.GetSno()
	}
	if options.eventDefinitionInstanceBuilder == nil {
		options.eventDefinitionInstanceBuilder = event.WrappingDefinitionInstanceBuilder
	}

	return options
}

type EngineOption func(*EngineOptions)

func WithIdGeneratorBuilder(idGeneratorBuilder id.IGeneratorBuilder) EngineOption {
	return func(o *EngineOptions) {
		o.idGeneratorBuilder = idGeneratorBuilder
	}
}

func WithEventDefinitionInstanceBuilder(eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder) EngineOption {
	return func(o *EngineOptions) {
		o.eventDefinitionInstanceBuilder = eventDefinitionInstanceBuilder
	}
}

type Engine struct {
	*EngineOptions
}

func NewEngine(opts ...EngineOption) *Engine {
	options := NewEngineOptions(opts...)

	engine := &Engine{
		EngineOptions: options,
	}

	return engine
}

func (p *Engine) NewProcess(definitions *schema.Definitions, opts ...Option) (process *Process, err error) {

	var processElement *schema.Process
	for _, element := range *definitions.Processes() {
		able, ok := element.IsExecutable()
		if !ok || !able {
			continue
		}

		processElement = &element
	}

	if processElement == nil {
		return nil, fmt.Errorf("no process found in definitions")
	}

	options := NewOptions(opts...)

	tracer := options.tracer
	if tracer == nil {
		tracer = tracing.NewTracer(p.ctx)
		opts = append(opts, WithTracer(tracer))
	}

	idGenerator := options.idGenerator
	if idGenerator == nil {
		idGenerator, err = p.idGeneratorBuilder.NewIdGenerator(p.ctx, tracer)
		if err != nil {
			return nil, fmt.Errorf("create id generator: %w", err)
		}
		opts = append(opts, WithIdGenerator(idGenerator))
	}

	process, err = NewProcess(processElement, definitions, opts...)
	if err != nil {
		return
	}

	return
}
