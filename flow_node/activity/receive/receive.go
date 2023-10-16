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

package receive

import (
	"context"
	"sync/atomic"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/schema"
)

type imessage interface {
	message()
}

type nextActionMessage struct {
	Headers    map[string]any
	Properties map[string]any
	response   chan flow_node.IAction
}

func (m nextActionMessage) message() {}

type cancelMessage struct {
	response chan bool
}

func (m cancelMessage) message() {}

type ReceiveTask struct {
	*flow_node.Wiring
	ctx           context.Context
	cancel        context.CancelFunc
	element       *schema.ReceiveTask
	runnerChannel chan imessage
}

func NewReceiveTask(ctx context.Context, task *schema.ReceiveTask) activity.Constructor {
	return func(wiring *flow_node.Wiring) (node activity.Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		done := &atomic.Bool{}
		done.Store(false)
		taskNode := &ReceiveTask{
			Wiring:        wiring,
			ctx:           ctx,
			cancel:        cancel,
			element:       task,
			runnerChannel: make(chan imessage, len(wiring.Incoming)*2+1),
		}
		go taskNode.runner(ctx)
		node = taskNode
		return
	}
}

func (node *ReceiveTask) runner(ctx context.Context) {
	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case cancelMessage:
				node.cancel()
				m.response <- true
			case nextActionMessage:
				go func() {
					aResponse := &flow_node.FlowActionResponse{
						Variables: map[string]data.IItem{},
					}
					action := flow_node.FlowAction{
						Response:      aResponse,
						SequenceFlows: flow_node.AllSequenceFlows(&node.Outgoing),
					}

					response := make(chan doResponse, 1)
					at := &ActiveTrace{
						Context:  node.ctx,
						Activity: node,
						Headers:  m.Headers,
						response: response,
					}

					if node.element.ExtensionElementsField != nil &&
						node.element.ExtensionElementsField.TaskDefinitionField != nil {
						at.Type = node.element.ExtensionElementsField.TaskDefinitionField.Type
					}

					node.Tracer.Trace(at)
					select {
					case <-ctx.Done():
						node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
						return
					case out := <-response:
						if out.err != nil {
							aResponse.Err = out.err
						}
						for key, value := range out.properties {
							aResponse.Variables[key] = value
						}
					}

					m.response <- action
				}()
			default:
			}
		case <-ctx.Done():
			node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
			return
		}
	}
}

func (node *ReceiveTask) NextAction(t flow_interface.T) chan flow_node.IAction {
	response := make(chan flow_node.IAction, 1)

	msg := nextActionMessage{
		response: response,
	}

	headers, _, _ := activity.FetchTaskDataInput(node.Locator, node.element)
	msg.Headers = headers

	node.runnerChannel <- msg
	return response
}

func (node *ReceiveTask) Element() schema.FlowNodeInterface {
	return node.element
}

func (node *ReceiveTask) Cancel() <-chan bool {
	response := make(chan bool)
	node.runnerChannel <- cancelMessage{response: response}
	return response
}
