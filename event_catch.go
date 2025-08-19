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
	"github.com/olive-io/bpmn/v2/pkg/logic"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type processEventMessage struct {
	event event.IEvent
}

func (m processEventMessage) message() {}

type catchEvent struct {
	*wiring
	element         *schema.CatchEvent
	mch             chan imessage
	activated       bool
	awaitingActions []chan IAction
	satisfier       *logic.CatchEventSatisfier
}

func newCatchEvent(wiring *wiring, element *schema.CatchEvent) (evt *catchEvent, err error) {
	evt = &catchEvent{
		wiring:          wiring,
		element:         element,
		mch:             make(chan imessage, len(wiring.incoming)*2+1),
		activated:       false,
		awaitingActions: make([]chan IAction, 0),
		satisfier:       logic.NewCatchEventSatisfier(element, wiring.eventDefinitionInstanceBuilder),
	}

	err = evt.eventEgress.RegisterEventConsumer(evt)
	if err != nil {
		return
	}
	return
}

func (evt *catchEvent) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-evt.mch:
			switch m := msg.(type) {
			case processEventMessage:
				if evt.activated {
					evt.tracer.Send(EventObservedTrace{Node: evt.element, Event: m.event})
					if satisfied, _ := evt.satisfier.Satisfy(m.event); satisfied {
						awaitingActions := evt.awaitingActions
						for _, actionChan := range awaitingActions {
							actionChan <- FlowAction{SequenceFlows: allSequenceFlows(&evt.outgoing)}
						}
						evt.awaitingActions = make([]chan IAction, 0)
						evt.activated = false
					}
				}
			case nextActionMessage:
				if !evt.activated {
					evt.activated = true
					evt.tracer.Send(ActiveListeningTrace{Node: evt.element})
				}
				evt.awaitingActions = append(evt.awaitingActions, m.response)
			default:
			}
		case <-ctx.Done():
			evt.tracer.Send(CancellationFlowNodeTrace{Node: evt.element})
			return
		}
	}
}

func (evt *catchEvent) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	evt.mch <- processEventMessage{event: ev}
	result = event.Consumed
	return
}

func (evt *catchEvent) NextAction(ctx context.Context, flow Flow) chan IAction {
	sender := evt.tracer.RegisterSender()
	go evt.run(ctx, sender)

	response := make(chan IAction)
	evt.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (evt *catchEvent) Element() schema.FlowNodeInterface {
	return evt.element
}

type ActiveListeningTrace struct {
	Node *schema.CatchEvent
}

func (t ActiveListeningTrace) Unpack() any { return t.Node }

// EventObservedTrace signals the fact that a particular event
// has in fact observed by the node
type EventObservedTrace struct {
	Node  *schema.CatchEvent
	Event event.IEvent
}

func (t EventObservedTrace) Unpack() any { return t.Node }
