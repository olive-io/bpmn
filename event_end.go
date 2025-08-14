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
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type EndEvent struct {
	*Wiring
	ctx                  context.Context
	element              *schema.EndEvent
	activated            bool
	completed            bool
	mch                  chan imessage
	startEventsActivated []*schema.StartEvent
}

func NewEndEvent(ctx context.Context, wiring *Wiring, endEvent *schema.EndEvent) (evt *EndEvent, err error) {
	evt = &EndEvent{
		Wiring:               wiring,
		ctx:                  ctx,
		element:              endEvent,
		activated:            false,
		completed:            false,
		mch:                  make(chan imessage, len(wiring.Incoming)*2+1),
		startEventsActivated: make([]*schema.StartEvent, 0),
	}
	return
}

func (evt *EndEvent) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-evt.mch:
			switch m := msg.(type) {
			case nextActionMessage:
				if !evt.activated {
					evt.activated = true
				}
				// If the node already completed, then we essentially fuse it
				if evt.completed {
					m.response <- CompleteAction{}
					continue
				}

				if _, err := evt.EventIngress.ConsumeEvent(event.MakeEndEvent(evt.element)); err == nil {
					evt.completed = true
					m.response <- CompleteAction{}
				} else {
					evt.Wiring.Tracer.Send(ErrorTrace{Error: err})
				}
			default:
			}
		case <-ctx.Done():
			evt.Tracer.Send(CancellationFlowNodeTrace{Node: evt.element})
			return
		}
	}
}

func (evt *EndEvent) NextAction(Flow) chan IAction {
	sender := evt.Tracer.RegisterSender()
	go evt.run(evt.ctx, sender)

	response := make(chan IAction, 1)
	evt.mch <- nextActionMessage{response: response}
	return response
}

func (evt *EndEvent) Element() schema.FlowNodeInterface {
	return evt.element
}
