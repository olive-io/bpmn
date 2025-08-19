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

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/logic"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type startMessage struct{}

func (m startMessage) message() {}

type eventMessage struct {
	event event.IEvent
}

func (m eventMessage) message() {}

type StartEvent struct {
	*Wiring
	element     *schema.StartEvent
	mch         chan imessage
	activated   bool
	idGenerator id.IGenerator
	satisfier   *logic.CatchEventSatisfier
}

func NewStartEvent(wiring *Wiring,
	startEvent *schema.StartEvent, idGenerator id.IGenerator,
) (evt *StartEvent, err error) {
	eventDefinitions := startEvent.EventDefinitions()
	eventInstances := make([]event.IDefinitionInstance, len(eventDefinitions))

	for i, eventDefinition := range eventDefinitions {
		var instance event.IDefinitionInstance
		instance, err = wiring.EventDefinitionInstanceBuilder.NewEventDefinitionInstance(eventDefinition)
		if err != nil {
			return
		}
		eventInstances[i] = instance
	}

	evt = &StartEvent{
		Wiring:      wiring,
		element:     startEvent,
		mch:         make(chan imessage, len(wiring.Incoming)*2+1),
		activated:   false,
		idGenerator: idGenerator,
		satisfier:   logic.NewCatchEventSatisfier(startEvent, wiring.EventDefinitionInstanceBuilder),
	}
	err = evt.EventEgress.RegisterEventConsumer(evt)
	if err != nil {
		return
	}
	return
}

func (evt *StartEvent) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-evt.mch:
			switch m := msg.(type) {
			case nextActionMessage:
				if !evt.activated {
					evt.activated = true
					m.response <- FlowAction{SequenceFlows: AllSequenceFlows(&evt.Outgoing)}
				} else {
					m.response <- CompleteAction{}
				}
			case startMessage:
				evt.flow(ctx)
			case eventMessage:
				if !evt.activated {
					if satisfied, _ := evt.satisfier.Satisfy(m.event); satisfied {
						evt.flow(ctx)
					}
				}
			default:
			}
		case <-ctx.Done():
			evt.Tracer.Send(CancellationFlowNodeTrace{Node: evt.element})
			return
		}
	}
}

func (evt *StartEvent) flow(ctx context.Context) {
	flowable := newFlow(evt.Wiring.Definitions, evt, evt.Wiring.Tracer,
		evt.Wiring.FlowNodeMapping, evt.Wiring.FlowWaitGroup, evt.idGenerator, nil, evt.Locator)
	flowable.Start(ctx)
}

func (evt *StartEvent) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	evt.mch <- eventMessage{event: ev}
	result = event.Consumed
	return
}

func (evt *StartEvent) Trigger(ctx context.Context) {
	sender := evt.Tracer.RegisterSender()
	go evt.run(ctx, sender)

	evt.mch <- startMessage{}
}

func (evt *StartEvent) NextAction(ctx context.Context, flow Flow) chan IAction {
	response := make(chan IAction)
	evt.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (evt *StartEvent) Element() schema.FlowNodeInterface {
	return evt.element
}
