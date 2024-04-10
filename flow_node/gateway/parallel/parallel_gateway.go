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

package parallel

import (
	"context"

	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/gateway"
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
	element               *schema.ParallelGateway
	runnerChannel         chan imessage
	reportedIncomingFlows int
	awaitingActions       []chan flow_node.IAction
	noOfIncomingFlows     int
}

func New(ctx context.Context, wiring *flow_node.Wiring, parallelGateway *schema.ParallelGateway) (node *Node, err error) {
	node = &Node{
		Wiring:                wiring,
		element:               parallelGateway,
		runnerChannel:         make(chan imessage, len(wiring.Incoming)*2+1),
		reportedIncomingFlows: 0,
		awaitingActions:       make([]chan flow_node.IAction, 0),
		noOfIncomingFlows:     len(wiring.Incoming),
	}
	sender := node.Tracer.RegisterSender()
	go node.runner(ctx, sender)
	return
}

func (node *Node) flowWhenReady() {
	if node.reportedIncomingFlows == node.noOfIncomingFlows {
		node.reportedIncomingFlows = 0
		awaitingActions := node.awaitingActions
		node.awaitingActions = make([]chan flow_node.IAction, 0)
		sequenceFlows := flow_node.AllSequenceFlows(&node.Outgoing)
		gateway.DistributeFlows(awaitingActions, sequenceFlows)
	}
}

func (node *Node) runner(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case nextActionMessage:
				node.reportedIncomingFlows++
				node.awaitingActions = append(node.awaitingActions, m.response)
				node.flowWhenReady()
				node.Tracer.Trace(IncomingFlowProcessedTrace{Node: node.element, Flow: m.flow})
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
