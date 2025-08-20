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

type startEvent struct {
	*wiring
	element     *schema.StartEvent
	mch         chan imessage
	activated   bool
	idGenerator id.IGenerator
	satisfier   *logic.CatchEventSatisfier
}

func newStartEvent(wr *wiring, element *schema.StartEvent, idGenerator id.IGenerator) (evt *startEvent, err error) {
	eventDefinitions := element.EventDefinitions()
	eventInstances := make([]event.IDefinitionInstance, len(eventDefinitions))

	for i, eventDefinition := range eventDefinitions {
		var instance event.IDefinitionInstance
		instance, err = wr.eventDefinitionInstanceBuilder.NewEventDefinitionInstance(eventDefinition)
		if err != nil {
			return
		}
		eventInstances[i] = instance
	}

	evt = &startEvent{
		wiring:      wr,
		element:     element,
		mch:         make(chan imessage, len(wr.incoming)*2+1),
		activated:   false,
		idGenerator: idGenerator,
		satisfier:   logic.NewCatchEventSatisfier(element, wr.eventDefinitionInstanceBuilder),
	}
	err = evt.eventEgress.RegisterEventConsumer(evt)
	if err != nil {
		return
	}
	return
}

func (evt *startEvent) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-evt.mch:
			switch m := msg.(type) {
			case nextActionMessage:
				if !evt.activated {
					evt.activated = true
					m.response <- flowAction{sequenceFlows: allSequenceFlows(&evt.outgoing)}
				} else {
					m.response <- completeAction{}
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
			evt.tracer.Send(CancellationFlowNodeTrace{Node: evt.element})
			return
		}
	}
}

func (evt *startEvent) flow(ctx context.Context) {
	flowable := newFlow(evt.wiring.definitions, evt, evt.wiring.tracer,
		evt.wiring.flowNodeMapping, evt.wiring.flowWaitGroup, evt.idGenerator, nil, evt.locator)
	flowable.Start(ctx)
}

func (evt *startEvent) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	evt.mch <- eventMessage{event: ev}
	result = event.Consumed
	return
}

func (evt *startEvent) Trigger(ctx context.Context) {
	sender := evt.tracer.RegisterSender()
	go evt.run(ctx, sender)

	evt.mch <- startMessage{}
}

func (evt *startEvent) NextAction(ctx context.Context, flow Flow) chan IAction {
	response := make(chan IAction)
	evt.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (evt *startEvent) Element() schema.FlowNodeInterface {
	return evt.element
}
