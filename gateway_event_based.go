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

func NewEventBasedGateway(ctx context.Context, wiring *Wiring, eventBasedGateway *schema.EventBasedGateway) (gw *EventBasedGateway, err error) {
	gw = &EventBasedGateway{
		Wiring:    wiring,
		element:   eventBasedGateway,
		mch:       make(chan imessage, len(wiring.Incoming)*2+1),
		activated: false,
	}
	sender := gw.Tracer.RegisterSender()
	go gw.run(ctx, sender)
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
				sequenceFlows := AllSequenceFlows(&gw.Outgoing)
				terminationChannels := make(map[schema.IdRef]chan bool)
				for _, sequenceFlow := range sequenceFlows {
					if idPtr, present := sequenceFlow.Id(); present {
						terminationChannels[*idPtr] = make(chan bool)
					} else {
						err := errors.NotFoundError{
							Expected: fmt.Sprintf("id for %#v", sequenceFlow),
						}
						gw.Tracer.Trace(ErrorTrace{Error: err})
					}
				}

				action := FlowAction{
					Terminate: func(sequenceFlowId *schema.IdRef) chan bool {
						return terminationChannels[*sequenceFlowId]
					},
					SequenceFlows: sequenceFlows,
					ActionTransformer: func(sequenceFlowId *schema.IdRef, action IAction) IAction {
						// only first one is to flow
						if atomic.CompareAndSwapInt32(&first, 0, 1) {
							gw.Tracer.Trace(DeterminationMadeTrace{Node: gw.element})
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
			gw.Tracer.Trace(CancellationFlowNodeTrace{Node: gw.element})
			return
		}
	}
}

func (gw *EventBasedGateway) NextAction(flow Flow) chan IAction {
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

func (t DeterminationMadeTrace) Element() any { return t.Node }
