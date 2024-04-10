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

package catch

import (
	"context"

	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/pkg/logic"
	"github.com/olive-io/bpmn/tracing"
)

type imessage interface {
	message()
}

type nextActionMessage struct {
	response chan flow_node.IAction
}

func (m nextActionMessage) message() {}

type processEventMessage struct {
	event event.IEvent
}

func (m processEventMessage) message() {}

type Node struct {
	*flow_node.Wiring
	element         *schema.CatchEvent
	runnerChannel   chan imessage
	activated       bool
	awaitingActions []chan flow_node.IAction
	satisfier       *logic.CatchEventSatisfier
}

func New(ctx context.Context, wiring *flow_node.Wiring, catchEvent *schema.CatchEvent) (node *Node, err error) {
	node = &Node{
		Wiring:          wiring,
		element:         catchEvent,
		runnerChannel:   make(chan imessage, len(wiring.Incoming)*2+1),
		activated:       false,
		awaitingActions: make([]chan flow_node.IAction, 0),
		satisfier:       logic.NewCatchEventSatisfier(catchEvent, wiring.EventDefinitionInstanceBuilder),
	}
	sender := node.Tracer.RegisterSender()
	go node.runner(ctx, sender)
	err = node.EventEgress.RegisterEventConsumer(node)
	if err != nil {
		return
	}
	return
}

func (node *Node) runner(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case processEventMessage:
				if node.activated {
					node.Tracer.Trace(EventObservedTrace{Node: node.element, Event: m.event})
					if satisfied, _ := node.satisfier.Satisfy(m.event); satisfied {
						awaitingActions := node.awaitingActions
						for _, actionChan := range awaitingActions {
							actionChan <- flow_node.FlowAction{SequenceFlows: flow_node.AllSequenceFlows(&node.Outgoing)}
						}
						node.awaitingActions = make([]chan flow_node.IAction, 0)
						node.activated = false
					}
				}
			case nextActionMessage:
				if !node.activated {
					node.activated = true
					node.Tracer.Trace(ActiveListeningTrace{Node: node.element})
				}
				node.awaitingActions = append(node.awaitingActions, m.response)
			default:
			}
		case <-ctx.Done():
			node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
			return
		}
	}
}

func (node *Node) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	node.runnerChannel <- processEventMessage{event: ev}
	result = event.Consumed
	return
}

func (node *Node) NextAction(flow_interface.T) chan flow_node.IAction {
	response := make(chan flow_node.IAction)
	node.runnerChannel <- nextActionMessage{response: response}
	return response
}

func (node *Node) Element() schema.FlowNodeInterface {
	return node.element
}
