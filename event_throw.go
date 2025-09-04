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
	"sync"
	"sync/atomic"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type throwEvent struct {
	*wiring
	element         *schema.ThrowEvent
	mch             chan imessage
	activated       atomic.Bool
	awaitingActions []chan IAction
	once            sync.Once
}

func newThrowEvent(wiring *wiring, element *schema.ThrowEvent) (evt *throwEvent, err error) {
	evt = &throwEvent{
		wiring:          wiring,
		element:         element,
		mch:             make(chan imessage, len(wiring.incoming)*2+1),
		activated:       atomic.Bool{},
		awaitingActions: make([]chan IAction, 0),
	}

	err = evt.eventEgress.RegisterEventConsumer(evt)
	if err != nil {
		return
	}
	return
}

func (evt *throwEvent) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-evt.mch:
			switch m := msg.(type) {
			case processEventMessage:
			case nextActionMessage:
				if !evt.activated.Load() {
					evt.activated.Store(true)
					m.response <- flowAction{sequenceFlows: allSequenceFlows(&evt.outgoing)}
				} else {
					m.response <- completeAction{}
				}
			default:
			}
		case <-ctx.Done():
			evt.tracer.Send(CancellationFlowNodeTrace{Node: evt.element})
			return
		}
	}
}

func (evt *throwEvent) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	evt.mch <- processEventMessage{event: ev}
	result = event.Consumed
	return
}

func (evt *throwEvent) NextAction(ctx context.Context, flow Flow) chan IAction {
	evt.once.Do(func() {
		sender := evt.tracer.RegisterSender()
		go evt.run(ctx, sender)
	})

	response := make(chan IAction)
	evt.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (evt *throwEvent) Element() schema.FlowNodeInterface {
	return evt.element
}
