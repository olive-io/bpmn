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

package end

import (
	"context"

	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/tracing"
)

type imessage interface {
	message()
}

type nextActionMessage struct {
	response chan flow_node.IAction
}

func (m nextActionMessage) message() {}

type Node struct {
	*flow_node.Wiring
	element              *schema.EndEvent
	activated            bool
	completed            bool
	runnerChannel        chan imessage
	startEventsActivated []*schema.StartEvent
}

func New(ctx context.Context, wiring *flow_node.Wiring, endEvent *schema.EndEvent) (node *Node, err error) {
	node = &Node{
		Wiring:               wiring,
		element:              endEvent,
		activated:            false,
		completed:            false,
		runnerChannel:        make(chan imessage, len(wiring.Incoming)*2+1),
		startEventsActivated: make([]*schema.StartEvent, 0),
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
				if !node.activated {
					node.activated = true
				}
				// If the node already completed, then we essentially fuse it
				if node.completed {
					m.response <- flow_node.CompleteAction{}
					continue
				}

				if _, err := node.EventIngress.ConsumeEvent(
					event.MakeEndEvent(node.element),
				); err == nil {
					node.completed = true
					m.response <- flow_node.CompleteAction{}
				} else {
					node.Wiring.Tracer.Trace(tracing.ErrorTrace{Error: err})
				}
			default:
			}
		case <-ctx.Done():
			node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
			return
		}
	}
}

func (node *Node) NextAction(flow_interface.T) chan flow_node.IAction {
	response := make(chan flow_node.IAction, 1)
	node.runnerChannel <- nextActionMessage{response: response}
	return response
}

func (node *Node) Element() schema.FlowNodeInterface {
	return node.element
}
