/*
   Copyright 2023 The bpmn Authors

   This program is offered under a commercial and under the AGPL license.
   For AGPL licensing, see below.

   AGPL licensing:
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package inclusive

import (
	"context"
	"fmt"

	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/gateway"
	"github.com/olive-io/bpmn/pkg/id"
	"github.com/olive-io/bpmn/sequence_flow"
	"github.com/olive-io/bpmn/tracing"
)

type NoEffectiveSequenceFlows struct {
	*schema.InclusiveGateway
}

func (e NoEffectiveSequenceFlows) Error() string {
	ownId := "<unnamed>"
	if ownIdPtr, present := e.InclusiveGateway.Id(); present {
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

type flowSync struct {
	response chan flow_node.IAction
	flow     flow_interface.T
}

type Node struct {
	*flow_node.Wiring
	element                 *schema.InclusiveGateway
	runnerChannel           chan imessage
	defaultSequenceFlow     *sequence_flow.SequenceFlow
	nonDefaultSequenceFlows []*sequence_flow.SequenceFlow
	probing                 *chan flow_node.IAction
	activated               *flowSync
	awaiting                []id.Id
	arrived                 []id.Id
	sync                    []chan flow_node.IAction
	*flowTracker
	synchronized bool
}

func New(ctx context.Context, wiring *flow_node.Wiring, inclusiveGateway *schema.InclusiveGateway) (node *Node, err error) {
	var defaultSequenceFlow *sequence_flow.SequenceFlow

	if seqFlow, present := inclusiveGateway.Default(); present {
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
		element:                 inclusiveGateway,
		runnerChannel:           make(chan imessage, len(wiring.Incoming)*2+1),
		nonDefaultSequenceFlows: nonDefaultSequenceFlows,
		defaultSequenceFlow:     defaultSequenceFlow,
		flowTracker:             newFlowTracker(ctx, wiring.Tracer, inclusiveGateway),
	}
	sender := node.Tracer.RegisterSender()
	go node.runner(ctx, sender)
	return
}

func (node *Node) runner(ctx context.Context, sender tracing.ISenderHandle) {
	defer node.flowTracker.shutdown()
	activity := node.flowTracker.activity()

	defer sender.Done()

	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case probingReport:
				response := node.probing
				if response == nil {
					// Reschedule, there's no next action yet
					go func() {
						node.runnerChannel <- m
					}()
					continue
				}
				node.probing = nil
				flow := make([]*sequence_flow.SequenceFlow, 0)
				for _, i := range m.result {
					flow = append(flow, node.nonDefaultSequenceFlows[i])
				}

				switch len(flow) {
				case 0:
					// no successful non-default sequence flows
					if node.defaultSequenceFlow == nil {
						// exception (Table 13.2)
						node.Wiring.Tracer.Trace(tracing.ErrorTrace{
							Error: NoEffectiveSequenceFlows{
								InclusiveGateway: node.element,
							},
						})
					} else {
						gateway.DistributeFlows(node.sync, []*sequence_flow.SequenceFlow{node.defaultSequenceFlow})
					}
				default:
					gateway.DistributeFlows(node.sync, flow)
				}
				node.synchronized = false
				node.activated = nil
			case nextActionMessage:
				if node.synchronized {
					if m.flow.Id() == node.activated.flow.Id() {
						// Activating flow returned
						node.sync = append(node.sync, m.response)
						node.probing = &m.response
						// and now we wait until the probe has returned
					}
				} else {
					if node.activated == nil {
						// Haven't been activated yet
						node.activated = &flowSync{response: m.response, flow: m.flow}
						node.awaiting = node.flowTracker.activeFlowsInCohort(m.flow.Id())
						node.arrived = []id.Id{m.flow.Id()}
						node.sync = make([]chan flow_node.IAction, 0)
					} else {
						// Already activated
						node.arrived = append(node.arrived, m.flow.Id())
						node.sync = append(node.sync, m.response)
					}
					node.trySync()
				}

			default:
			}
		case <-activity:
			if !node.synchronized && node.activated != nil {
				node.awaiting = node.flowTracker.activeFlowsInCohort(node.activated.flow.Id())
				node.trySync()
			}
		case <-ctx.Done():
			node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
			return
		}
	}
}

func (node *Node) trySync() {
	if !node.synchronized && len(node.arrived) >= len(node.awaiting) {
		// Have we got everybody?
		matches := 0
		for i := range node.arrived {
			for j := range node.awaiting {
				if node.awaiting[j] == node.arrived[i] {
					matches++
				}
			}
		}
		if matches == len(node.awaiting) {
			anId := node.activated.flow.Id()
			// Probe outgoing sequence flow using the first flow
			node.activated.response <- flow_node.ProbeAction{
				SequenceFlows: node.nonDefaultSequenceFlows,
				ProbeReport: func(indices []int) {
					node.runnerChannel <- probingReport{
						result: indices,
						flowId: anId,
					}
				},
			}

			node.synchronized = true
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
