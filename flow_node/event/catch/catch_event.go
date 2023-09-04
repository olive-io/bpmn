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

package catch

import (
	"context"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tools/logic"
	"github.com/olive-io/bpmn/tracing"
)

type message interface {
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
	runnerChannel   chan message
	activated       bool
	awaitingActions []chan flow_node.IAction
	satisfier       *logic.CatchEventSatisfier
}

func New(ctx context.Context, wiring *flow_node.Wiring, catchEvent *schema.CatchEvent) (node *Node, err error) {
	node = &Node{
		Wiring:          wiring,
		element:         catchEvent,
		runnerChannel:   make(chan message, len(wiring.Incoming)*2+1),
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
