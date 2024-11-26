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

type CatchEvent struct {
	*Wiring
	element         *schema.CatchEvent
	runnerChannel   chan imessage
	activated       bool
	awaitingActions []chan IAction
	satisfier       *logic.CatchEventSatisfier
}

func NewCatchEvent(ctx context.Context, wiring *Wiring, catchEvent *schema.CatchEvent) (evt *CatchEvent, err error) {
	evt = &CatchEvent{
		Wiring:          wiring,
		element:         catchEvent,
		runnerChannel:   make(chan imessage, len(wiring.Incoming)*2+1),
		activated:       false,
		awaitingActions: make([]chan IAction, 0),
		satisfier:       logic.NewCatchEventSatisfier(catchEvent, wiring.EventDefinitionInstanceBuilder),
	}
	sender := evt.Tracer.RegisterSender()
	go evt.runner(ctx, sender)
	err = evt.EventEgress.RegisterEventConsumer(evt)
	if err != nil {
		return
	}
	return
}

func (evt *CatchEvent) runner(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-evt.runnerChannel:
			switch m := msg.(type) {
			case processEventMessage:
				if evt.activated {
					evt.Tracer.Trace(EventObservedTrace{Node: evt.element, Event: m.event})
					if satisfied, _ := evt.satisfier.Satisfy(m.event); satisfied {
						awaitingActions := evt.awaitingActions
						for _, actionChan := range awaitingActions {
							actionChan <- FlowAction{SequenceFlows: AllSequenceFlows(&evt.Outgoing)}
						}
						evt.awaitingActions = make([]chan IAction, 0)
						evt.activated = false
					}
				}
			case nextActionMessage:
				if !evt.activated {
					evt.activated = true
					evt.Tracer.Trace(ActiveListeningTrace{Node: evt.element})
				}
				evt.awaitingActions = append(evt.awaitingActions, m.response)
			default:
			}
		case <-ctx.Done():
			evt.Tracer.Trace(CancellationFlowNodeTrace{Node: evt.element})
			return
		}
	}
}

func (evt *CatchEvent) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	evt.runnerChannel <- processEventMessage{event: ev}
	result = event.Consumed
	return
}

func (evt *CatchEvent) NextAction(flow T) chan IAction {
	response := make(chan IAction)
	evt.runnerChannel <- nextActionMessage{response: response, flow: flow}
	return response
}

func (evt *CatchEvent) Element() schema.FlowNodeInterface {
	return evt.element
}

type ActiveListeningTrace struct {
	Node *schema.CatchEvent
}

func (t ActiveListeningTrace) TraceInterface() {}

// EventObservedTrace signals the fact that a particular event
// has been in fact observed by the node
type EventObservedTrace struct {
	Node  *schema.CatchEvent
	Event event.IEvent
}

func (t EventObservedTrace) TraceInterface() {}
