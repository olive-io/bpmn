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

package timer

import (
	"context"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/clock"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type eventDefinitionInstanceBuilder struct {
	// context here keeps context from the creation of the instance builder
	// It is not a final decision, but it currently seems to make more sense
	// to "attach" it to this context instead of the context that can be passed
	// through event.IDefinitionInstance.NewEventDefinitionInstance. Time will tell.
	context      context.Context
	eventIngress event.IConsumer
	tracer       tracing.ITracer
}

type eventDefinitionInstance struct {
	definition schema.TimerEventDefinition
}

func (e *eventDefinitionInstance) EventDefinition() schema.EventDefinitionInterface {
	return &e.definition
}

func (e *eventDefinitionInstanceBuilder) NewEventDefinitionInstance(def schema.EventDefinitionInterface) (definitionInstance event.IDefinitionInstance, err error) {
	if timerEventDefinition, ok := def.(*schema.TimerEventDefinition); ok {
		var c clock.IClock
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
						e.tracer.Send(bpmn.ErrorTrace{Error: err})
					}
				}
			}
		}(e.context)
	}
	return
}

func EventDefinitionInstanceBuilder(
	ctx context.Context,
	eventIngress event.IConsumer,
	tracer tracing.ITracer,
) event.IDefinitionInstanceBuilder {
	return &eventDefinitionInstanceBuilder{
		context:      ctx,
		eventIngress: eventIngress,
		tracer:       tracer,
	}
}
