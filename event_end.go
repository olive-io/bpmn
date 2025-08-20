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

type endEvent struct {
	*wiring
	element              *schema.EndEvent
	activated            bool
	completed            bool
	mch                  chan imessage
	startEventsActivated []*schema.StartEvent
}

func newEndEvent(wr *wiring, element *schema.EndEvent) (evt *endEvent, err error) {
	evt = &endEvent{
		wiring:               wr,
		element:              element,
		activated:            false,
		completed:            false,
		mch:                  make(chan imessage, len(wr.incoming)*2+1),
		startEventsActivated: make([]*schema.StartEvent, 0),
	}
	return
}

func (evt *endEvent) run(ctx context.Context, sender tracing.ISenderHandle) {
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
					m.response <- completeAction{}
					continue
				}

				if _, err := evt.eventIngress.ConsumeEvent(event.MakeEndEvent(evt.element)); err == nil {
					evt.completed = true
					m.response <- completeAction{}
				} else {
					evt.wiring.tracer.Send(ErrorTrace{Error: err})
				}
			default:
			}
		case <-ctx.Done():
			evt.tracer.Send(CancellationFlowNodeTrace{Node: evt.element})
			return
		}
	}
}

func (evt *endEvent) NextAction(ctx context.Context, flow Flow) chan IAction {
	sender := evt.tracer.RegisterSender()
	go evt.run(ctx, sender)

	response := make(chan IAction, 1)
	evt.mch <- nextActionMessage{response: response}
	return response
}

func (evt *endEvent) Element() schema.FlowNodeInterface {
	return evt.element
}
