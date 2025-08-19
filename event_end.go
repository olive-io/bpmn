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
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type EndEvent struct {
	*Wiring
	element              *schema.EndEvent
	activated            bool
	completed            bool
	mch                  chan imessage
	startEventsActivated []*schema.StartEvent
}

func NewEndEvent(wiring *Wiring, endEvent *schema.EndEvent) (evt *EndEvent, err error) {
	evt = &EndEvent{
		Wiring:               wiring,
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

func (evt *EndEvent) NextAction(ctx context.Context, flow Flow) chan IAction {
	sender := evt.Tracer.RegisterSender()
	go evt.run(ctx, sender)

	response := make(chan IAction, 1)
	evt.mch <- nextActionMessage{response: response}
	return response
}

func (evt *EndEvent) Element() schema.FlowNodeInterface {
	return evt.element
}
