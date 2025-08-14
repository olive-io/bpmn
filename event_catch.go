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
	ctx             context.Context
	element         *schema.CatchEvent
	mch             chan imessage
	activated       bool
	awaitingActions []chan IAction
	satisfier       *logic.CatchEventSatisfier
}

func NewCatchEvent(ctx context.Context, wiring *Wiring, catchEvent *schema.CatchEvent) (evt *CatchEvent, err error) {
	evt = &CatchEvent{
		Wiring:          wiring,
		ctx:             ctx,
		element:         catchEvent,
		mch:             make(chan imessage, len(wiring.Incoming)*2+1),
		activated:       false,
		awaitingActions: make([]chan IAction, 0),
		satisfier:       logic.NewCatchEventSatisfier(catchEvent, wiring.EventDefinitionInstanceBuilder),
	}

	err = evt.EventEgress.RegisterEventConsumer(evt)
	if err != nil {
		return
	}
	return
}

func (evt *CatchEvent) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-evt.mch:
			switch m := msg.(type) {
			case processEventMessage:
				if evt.activated {
					evt.Tracer.Send(EventObservedTrace{Node: evt.element, Event: m.event})
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
					evt.Tracer.Send(ActiveListeningTrace{Node: evt.element})
				}
				evt.awaitingActions = append(evt.awaitingActions, m.response)
			default:
			}
		case <-ctx.Done():
			evt.Tracer.Send(CancellationFlowNodeTrace{Node: evt.element})
			return
		}
	}
}

func (evt *CatchEvent) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	evt.mch <- processEventMessage{event: ev}
	result = event.Consumed
	return
}

func (evt *CatchEvent) NextAction(flow Flow) chan IAction {
	sender := evt.Tracer.RegisterSender()
	go evt.run(evt.ctx, sender)

	response := make(chan IAction)
	evt.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (evt *CatchEvent) Element() schema.FlowNodeInterface {
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
