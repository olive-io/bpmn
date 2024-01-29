// Copyright 2023 The bpmn Authors
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

package exclusive

import (
	"context"
	"fmt"

	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/sequence_flow"
	"github.com/olive-io/bpmn/tools/id"
	"github.com/olive-io/bpmn/tracing"
)

type NoEffectiveSequenceFlows struct {
	*schema.ExclusiveGateway
}

func (e NoEffectiveSequenceFlows) Error() string {
	ownId := "<unnamed>"
	if ownIdPtr, present := e.ExclusiveGateway.Id(); present {
		ownId = *ownIdPtr
	}
	return fmt.Sprintf("No effective sequence flows found in exclusive gateway `%v`", ownId)
}

type imessage interface {
	message()
}

type nextActionMessage struct {
	response chan flow_node.IAction
	flow     flow_interface.T
}

func (m nextActionMessage) message() {}

type probingReport struct {
	result []int
	flowId id.Id
}

func (m probingReport) message() {}

type Node struct {
	*flow_node.Wiring
	element                 *schema.ExclusiveGateway
	runnerChannel           chan imessage
	defaultSequenceFlow     *sequence_flow.SequenceFlow
	nonDefaultSequenceFlows []*sequence_flow.SequenceFlow
	probing                 map[id.Id]*chan flow_node.IAction
}

func New(ctx context.Context, wiring *flow_node.Wiring, exclusiveGateway *schema.ExclusiveGateway) (node *Node, err error) {
	var defaultSequenceFlow *sequence_flow.SequenceFlow

	if seqFlow, present := exclusiveGateway.Default(); present {
		if node, found := wiring.Process.FindBy(schema.ExactId(*seqFlow).
			And(schema.ElementType((*schema.SequenceFlow)(nil)))); found {
			defaultSequenceFlow = new(sequence_flow.SequenceFlow)
			*defaultSequenceFlow = sequence_flow.Make(
				node.(*schema.SequenceFlow),
				wiring.Process,
			)
		} else {
			err = errors.NotFoundError{
				Expected: fmt.Sprintf("default sequence flow with ID %s", *seqFlow),
			}
			return nil, err
		}
	}

	nonDefaultSequenceFlows := flow_node.AllSequenceFlows(&wiring.Outgoing,
		func(sequenceFlow *sequence_flow.SequenceFlow) bool {
			if defaultSequenceFlow == nil {
				return false
			}
			return *sequenceFlow == *defaultSequenceFlow
		},
	)

	node = &Node{
		Wiring:                  wiring,
		element:                 exclusiveGateway,
		runnerChannel:           make(chan imessage, len(wiring.Incoming)*2+1),
		nonDefaultSequenceFlows: nonDefaultSequenceFlows,
		defaultSequenceFlow:     defaultSequenceFlow,
		probing:                 make(map[id.Id]*chan flow_node.IAction),
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
			case probingReport:
				if response, ok := node.probing[m.flowId]; ok {
					if response == nil {
						// Reschedule, there's no next action yet
						go func() {
							node.runnerChannel <- m
						}()
						continue
					}
					delete(node.probing, m.flowId)
					flow := make([]*sequence_flow.SequenceFlow, 0)
					for _, i := range m.result {
						flow = append(flow, node.nonDefaultSequenceFlows[i])
						break
					}
					switch len(flow) {
					case 0:
						// no successful non-default sequence flows
						if node.defaultSequenceFlow == nil {
							// exception (Table 13.2)
							node.Wiring.Tracer.Trace(tracing.ErrorTrace{
								Error: NoEffectiveSequenceFlows{
									ExclusiveGateway: node.element,
								},
							})
						} else {
							// default
							*response <- flow_node.FlowAction{
								SequenceFlows:      []*sequence_flow.SequenceFlow{node.defaultSequenceFlow},
								UnconditionalFlows: []int{0},
							}
						}
					case 1:
						*response <- flow_node.FlowAction{
							SequenceFlows:      flow,
							UnconditionalFlows: []int{0},
						}
					default:
						node.Wiring.Tracer.Trace(tracing.ErrorTrace{
							Error: errors.InvalidArgumentError{
								Expected: fmt.Sprintf("maximum 1 outgoing exclusive gateway (%s) flow",
									node.Wiring.FlowNodeId),
								Actual: len(flow),
							},
						})
					}
				} else {
					node.Wiring.Tracer.Trace(tracing.ErrorTrace{
						Error: errors.InvalidStateError{
							Expected: fmt.Sprintf("probing[%s] is to be present (exclusive gateway %s)",
								m.flowId.String(), node.Wiring.FlowNodeId),
						},
					})
				}
			case nextActionMessage:
				if _, ok := node.probing[m.flow.Id()]; ok {
					node.probing[m.flow.Id()] = &m.response
					// and now we wait until the probe has returned
				} else {
					node.probing[m.flow.Id()] = nil
					m.response <- flow_node.ProbeAction{
						SequenceFlows: node.nonDefaultSequenceFlows,
						ProbeReport: func(indices []int) {
							node.runnerChannel <- probingReport{
								result: indices,
								flowId: m.flow.Id(),
							}
						},
					}
				}
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
