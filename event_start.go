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
	ctx         context.Context
	element     *schema.StartEvent
	mch         chan imessage
	activated   bool
	idGenerator id.IGenerator
	satisfier   *logic.CatchEventSatisfier
}

func NewStartEvent(ctx context.Context, wiring *Wiring,
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
		ctx:         ctx,
		element:     startEvent,
		mch:         make(chan imessage, len(wiring.Incoming)*2+1),
		activated:   false,
		idGenerator: idGenerator,
		satisfier:   logic.NewCatchEventSatisfier(startEvent, wiring.EventDefinitionInstanceBuilder),
	}
	sender := evt.Tracer.RegisterSender()
	go evt.run(evt.ctx, sender)
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

func (evt *StartEvent) Trigger() {
	evt.mch <- startMessage{}
}

func (evt *StartEvent) NextAction(flow Flow) chan IAction {
	response := make(chan IAction)
	evt.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (evt *StartEvent) Element() schema.FlowNodeInterface {
	return evt.element
}
