// Copyright 2023 Lack (xingyys@gmail.com).
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

package timer

import (
	"context"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tools/clock"
	"github.com/olive-io/bpmn/tracing"
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
	eventIngress event.IConsumer,
	tracer tracing.ITracer,
) event.IDefinitionInstanceBuilder {
	return &eventDefinitionInstanceBuilder{
		context:      ctx,
		eventIngress: eventIngress,
		tracer:       tracer,
	}
}
