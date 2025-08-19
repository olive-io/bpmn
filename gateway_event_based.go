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
	"fmt"
	"sync/atomic"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type EventBasedGateway struct {
	*Wiring
	element   *schema.EventBasedGateway
	mch       chan imessage
	activated bool
}

func NewEventBasedGateway(wiring *Wiring, eventBasedGateway *schema.EventBasedGateway) (gw *EventBasedGateway, err error) {
	gw = &EventBasedGateway{
		Wiring:    wiring,
		element:   eventBasedGateway,
		mch:       make(chan imessage, len(wiring.Incoming)*2+1),
		activated: false,
	}
	return
}

func (gw *EventBasedGateway) run(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-gw.mch:
			switch m := msg.(type) {
			case nextActionMessage:
				var first int32 = 0
				sequences := AllSequenceFlows(&gw.Outgoing)
				terminationChannels := make(map[schema.IdRef]chan bool)
				for _, sequenceFlow := range sequences {
					if idPtr, present := sequenceFlow.Id(); present {
						terminationChannels[*idPtr] = make(chan bool)
					} else {
						err := errors.NotFoundError{
							Expected: fmt.Sprintf("id for %#v", sequenceFlow),
						}
						gw.Tracer.Send(ErrorTrace{Error: err})
					}
				}

				action := FlowAction{
					Terminate: func(sequenceFlowId *schema.IdRef) chan bool {
						return terminationChannels[*sequenceFlowId]
					},
					SequenceFlows: sequences,
					ActionTransformer: func(sequenceFlowId *schema.IdRef, action IAction) IAction {
						// only the first one is to flow
						if atomic.CompareAndSwapInt32(&first, 0, 1) {
							gw.Tracer.Send(DeterminationMadeTrace{Node: gw.element})
							for terminationCandidateId, ch := range terminationChannels {
								if sequenceFlowId != nil && terminationCandidateId != *sequenceFlowId {
									ch <- true
								}
								close(ch)
							}
							terminationChannels = make(map[schema.IdRef]chan bool)
							return action
						} else {
							return CompleteAction{}
						}
					},
				}

				m.response <- action
			default:
			}
		case <-ctx.Done():
			gw.Tracer.Send(CancellationFlowNodeTrace{Node: gw.element})
			return
		}
	}
}

func (gw *EventBasedGateway) NextAction(ctx context.Context, flow Flow) chan IAction {
	sender := gw.Tracer.RegisterSender()
	go gw.run(ctx, sender)

	response := make(chan IAction)
	gw.mch <- nextActionMessage{response: response, flow: flow}
	return response
}

func (gw *EventBasedGateway) Element() schema.FlowNodeInterface {
	return gw.element
}

type DeterminationMadeTrace struct {
	Node schema.FlowNodeInterface
}

func (t DeterminationMadeTrace) Unpack() any { return t.Node }
