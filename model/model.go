package model

import (
	"context"
	"sync"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/id"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/timer"
	"github.com/olive-io/bpmn/tracing"
)

type Model struct {
	Element                        *schema.Definitions
	processes                      []process.Process
	eventConsumersLock             sync.RWMutex
	eventConsumers                 []event.Consumer
	idGeneratorBuilder             id.GeneratorBuilder
	eventDefinitionInstanceBuilder event.DefinitionInstanceBuilder
	tracer                         tracing.Tracer
}

type Option func(context.Context, *Model) context.Context

func WithIdGenerator(builder id.GeneratorBuilder) Option {
	return func(ctx context.Context, model *Model) context.Context {
		model.idGeneratorBuilder = builder
		return ctx
	}
}

func WithEventDefinitionInstanceBuilder(builder event.DefinitionInstanceBuilder) Option {
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
func WithTracer(tracer tracing.Tracer) Option {
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

	model.processes = make([]process.Process, len(*procs))

	for i := range *procs {
		model.processes[i] = process.Make(&(*procs)[i], element,
			process.WithIdGenerator(model.idGeneratorBuilder),
			process.WithEventIngress(model), process.WithEventEgress(model),
			process.WithEventDefinitionInstanceBuilder(model),
			process.WithContext(ctx),
			process.WithTracer(model.tracer),
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

func (model *Model) FindProcessBy(f func(*process.Process) bool) (result *process.Process, found bool) {
	for i := range model.processes {
		if f(&model.processes[i]) {
			result = &model.processes[i]
			found = true
			return
		}
	}
	return
}

func (model *Model) ConsumeEvent(ev event.Event) (result event.ConsumptionResult, err error) {
	model.eventConsumersLock.RLock()
	// We're copying the list of consumers here to ensure that
	// new consumers can subscribe during event forwarding
	eventConsumers := model.eventConsumers
	model.eventConsumersLock.RUnlock()
	result, err = event.ForwardEvent(ev, &eventConsumers)
	return
}

func (model *Model) RegisterEventConsumer(ev event.Consumer) (err error) {
	model.eventConsumersLock.Lock()
	defer model.eventConsumersLock.Unlock()
	model.eventConsumers = append(model.eventConsumers, ev)
	return
}

func (model *Model) NewEventDefinitionInstance(def schema.EventDefinitionInterface) (event.DefinitionInstance, error) {
	if model.eventDefinitionInstanceBuilder != nil {
		return model.eventDefinitionInstanceBuilder.NewEventDefinitionInstance(def)
	} else {
		return event.WrapEventDefinition(def), nil
	}
}
