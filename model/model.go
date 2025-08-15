/*
Copyright 2024 The bpmn Authors

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

package model

import (
	"context"
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/timer"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type Model struct {
	Element                        *schema.Definitions
	processes                      []bpmn.Process
	eventConsumersLock             sync.RWMutex
	eventConsumers                 []event.IConsumer
	idGeneratorBuilder             id.IGeneratorBuilder
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder
	tracer                         tracing.ITracer
}

type Option func(context.Context, *Model) context.Context

func WithIdGenerator(builder id.IGeneratorBuilder) Option {
	return func(ctx context.Context, model *Model) context.Context {
		model.idGeneratorBuilder = builder
		return ctx
	}
}

func WithEventDefinitionInstanceBuilder(builder event.IDefinitionInstanceBuilder) Option {
	return func(ctx context.Context, model *Model) context.Context {
		model.eventDefinitionInstanceBuilder = builder
		return ctx
	}
}

// WithContext will pass a given context to a new model
// instead of implicitly generated one
func WithContext(newCtx context.Context) Option {
	return func(ctx context.Context, model *Model) context.Context {
		return newCtx
	}
}

// WithTracer overrides model's tracer
func WithTracer(tracer tracing.ITracer) Option {
	return func(ctx context.Context, model *Model) context.Context {
		model.tracer = tracer
		return ctx
	}
}

func New(element *schema.Definitions, options ...Option) *Model {
	procs := element.Processes()
	model := &Model{
		Element: element,
	}

	ctx := context.Background()

	for _, option := range options {
		ctx = option(ctx, model)
	}

	if model.idGeneratorBuilder == nil {
		model.idGeneratorBuilder = id.DefaultIdGeneratorBuilder
	}

	if model.tracer == nil {
		model.tracer = tracing.NewTracer(ctx)
	}

	if model.eventDefinitionInstanceBuilder == nil {
		model.eventDefinitionInstanceBuilder = event.DefinitionInstanceBuildingChain(
			timer.EventDefinitionInstanceBuilder(ctx, model, model.tracer),
			event.WrappingDefinitionInstanceBuilder,
		)
	}

	model.processes = make([]bpmn.Process, len(*procs))

	for i := range *procs {
		model.processes[i] = bpmn.MakeProcess(&(*procs)[i], element,
			bpmn.WithIdGenerator(model.idGeneratorBuilder),
			bpmn.WithEventIngress(model), bpmn.WithEventEgress(model),
			bpmn.WithEventDefinitionInstanceBuilder(model),
			bpmn.WithContext(ctx),
			bpmn.WithTracer(model.tracer),
		)
	}
	return model
}

func (model *Model) Run(ctx context.Context) (err error) {
	// Setup process instantiation
	for i := range *model.Element.Processes() {
		instantiatingFlowNodes := (*model.Element.Processes())[i].InstantiatingFlowNodes()
		for j := range instantiatingFlowNodes {
			flowNode := instantiatingFlowNodes[j]

			switch node := flowNode.(type) {
			case *schema.StartEvent:
				err = model.RegisterEventConsumer(newStartEventConsumer(ctx,
					model.tracer,
					&model.processes[i],
					node, model.eventDefinitionInstanceBuilder))
				if err != nil {
					return
				}
			case *schema.EventBasedGateway:
			case *schema.ReceiveTask:
			}
		}
	}
	return
}

func (model *Model) FindProcessBy(f func(*bpmn.Process) bool) (result *bpmn.Process, found bool) {
	for i := range model.processes {
		if f(&model.processes[i]) {
			result = &model.processes[i]
			found = true
			return
		}
	}
	return
}

func (model *Model) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	model.eventConsumersLock.RLock()
	// We're copying the list of consumers here to ensure that
	// new consumers can subscribe during event forwarding
	eventConsumers := model.eventConsumers
	model.eventConsumersLock.RUnlock()
	result, err = event.ForwardEvent(ev, &eventConsumers)
	return
}

func (model *Model) RegisterEventConsumer(ev event.IConsumer) (err error) {
	model.eventConsumersLock.Lock()
	defer model.eventConsumersLock.Unlock()
	model.eventConsumers = append(model.eventConsumers, ev)
	return
}

func (model *Model) NewEventDefinitionInstance(def schema.EventDefinitionInterface) (event.IDefinitionInstance, error) {
	if model.eventDefinitionInstanceBuilder != nil {
		return model.eventDefinitionInstanceBuilder.NewEventDefinitionInstance(def)
	} else {
		return event.WrapEventDefinition(def), nil
	}
}
