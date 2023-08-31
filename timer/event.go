package timer

import (
	"context"

	"github.com/olive-io/bpmn/clock"
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
)

type eventDefinitionInstanceBuilder struct {
	// context here keeps context from the creation of the instance builder
	// It is not a final decision, but it currently seems to make more sense
	// to "attach" it to this context instead of the context that can be passed
	// through event.DefinitionInstance.NewEventDefinitionInstance. Time will tell.
	context      context.Context
	eventIngress event.Consumer
	tracer       tracing.Tracer
}

type eventDefinitionInstance struct {
	definition schema.TimerEventDefinition
}

func (e *eventDefinitionInstance) EventDefinition() schema.EventDefinitionInterface {
	return &e.definition
}

func (e *eventDefinitionInstanceBuilder) NewEventDefinitionInstance(def schema.EventDefinitionInterface) (definitionInstance event.DefinitionInstance, err error) {
	if timerEventDefinition, ok := def.(*schema.TimerEventDefinition); ok {
		var c clock.Clock
		c, err = clock.FromContext(e.context)
		if err != nil {
			return
		}
		var timer chan schema.TimerEventDefinition
		timer, err = New(e.context, c, *timerEventDefinition)
		if err != nil {
			return
		}
		definitionInstance = &eventDefinitionInstance{*timerEventDefinition}
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case _, ok := <-timer:
					if !ok {
						return
					}
					_, err := e.eventIngress.ConsumeEvent(event.MakeTimerEvent(definitionInstance))
					if err != nil {
						e.tracer.Trace(tracing.ErrorTrace{Error: err})
					}
				}
			}
		}(e.context)
	}
	return
}

func EventDefinitionInstanceBuilder(
	ctx context.Context,
	eventIngress event.Consumer,
	tracer tracing.Tracer,
) event.DefinitionInstanceBuilder {
	return &eventDefinitionInstanceBuilder{
		context:      ctx,
		eventIngress: eventIngress,
		tracer:       tracer,
	}
}
