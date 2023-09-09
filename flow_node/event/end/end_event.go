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

package end

import (
	"context"

	"github.com/olive-io/bpmn/event"
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
