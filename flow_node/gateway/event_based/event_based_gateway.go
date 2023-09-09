// Copyright 2023 Lack (xingyys@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event_based

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
)

type imessage interface {
	message()
}

type nextActionMessage struct {
	response chan flow_node.IAction
	flow     flow_interface.T
}

func (m nextActionMessage) message() {}

type Node struct {
	*flow_node.Wiring
	element       *schema.EventBasedGateway
	runnerChannel chan imessage
	activated     bool
}

func New(ctx context.Context, wiring *flow_node.Wiring, eventBasedGateway *schema.EventBasedGateway) (node *Node, err error) {
	node = &Node{
		Wiring:        wiring,
		element:       eventBasedGateway,
		runnerChannel: make(chan imessage, len(wiring.Incoming)*2+1),
		activated:     false,
	}
	sender := node.Tracer.RegisterSender()
	go node.runner(ctx, sender)
	return
}

func (node *Node) runner(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case nextActionMessage:
				var first int32 = 0
				sequenceFlows := flow_node.AllSequenceFlows(&node.Outgoing)
				terminationChannels := make(map[schema.IdRef]chan bool)
				for _, sequenceFlow := range sequenceFlows {
					if idPtr, present := sequenceFlow.Id(); present {
						terminationChannels[*idPtr] = make(chan bool)
					} else {
						node.Tracer.Trace(tracing.ErrorTrace{Error: errors.NotFoundError{
							Expected: fmt.Sprintf("id for %#v", sequenceFlow),
						}})
					}
				}

				action := flow_node.FlowAction{
					Terminate: func(sequenceFlowId *schema.IdRef) chan bool {
						return terminationChannels[*sequenceFlowId]
					},
					SequenceFlows: sequenceFlows,
					ActionTransformer: func(sequenceFlowId *schema.IdRef, action flow_node.IAction) flow_node.IAction {
						// only first one is to flow
						if atomic.CompareAndSwapInt32(&first, 0, 1) {
							node.Tracer.Trace(DeterminationMadeTrace{Element: node.element})
							for terminationCandidateId, ch := range terminationChannels {
								if sequenceFlowId != nil && terminationCandidateId != *sequenceFlowId {
									ch <- true
								}
								close(ch)
							}
							terminationChannels = make(map[schema.IdRef]chan bool)
							return action
						} else {
							return flow_node.CompleteAction{}
						}
					},
				}

				m.response <- action
			default:
			}
		case <-ctx.Done():
			node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
			return
		}
	}
}

func (node *Node) NextAction(flow flow_interface.T) chan flow_node.IAction {
	response := make(chan flow_node.IAction)
	node.runnerChannel <- nextActionMessage{response: response, flow: flow}
	return response
}

func (node *Node) Element() schema.FlowNodeInterface {
	return node.element
}
