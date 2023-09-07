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

package script_task

import (
	"context"
	"strings"
	"sync/atomic"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/schema"
)

type message interface {
	message()
}

type nextActionMessage struct {
	DataObjects map[string]any
	Headers     map[string]any
	Properties  map[string]any
	response    chan flow_node.IAction
}

func (m nextActionMessage) message() {}

type cancelMessage struct {
	response chan bool
}

func (m cancelMessage) message() {}

type ScriptTask struct {
	*flow_node.Wiring
	ctx           context.Context
	cancel        context.CancelFunc
	element       *schema.ScriptTask
	runnerChannel chan message
}

func NewScriptTask(ctx context.Context, task *schema.ScriptTask) activity.Constructor {
	return func(wiring *flow_node.Wiring) (node activity.Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		done := &atomic.Bool{}
		done.Store(false)
		taskNode := &ScriptTask{
			Wiring:        wiring,
			ctx:           ctx,
			cancel:        cancel,
			element:       task,
			runnerChannel: make(chan message, len(wiring.Incoming)*2+1),
		}
		go taskNode.runner(ctx)
		node = taskNode
		return
	}
}

func (node *ScriptTask) runner(ctx context.Context) {
	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case cancelMessage:
				node.cancel()
				m.response <- true
			case nextActionMessage:
				go func() {
					action := flow_node.FlowAction{
						Variables:     map[string]data.IItem{},
						SequenceFlows: flow_node.AllSequenceFlows(&node.Outgoing),
					}

					response := make(chan submitResponse, 1)
					at := &ActiveTrace{
						Context:     node.ctx,
						Activity:    node,
						DataObjects: m.DataObjects,
						Headers:     m.Headers,
						Properties:  m.Properties,
						response:    response,
					}

					node.Tracer.Trace(at)
					select {
					case <-ctx.Done():
						return
					case out := <-response:
						if out.err != nil {
							action.Err = out.err
						}
						for key, value := range out.result {
							action.Variables[key] = value
						}
						m.response <- action
					}
				}()
			default:
			}
		case <-ctx.Done():
			return
		}
	}
}

func (node *ScriptTask) NextAction(t flow_interface.T) chan flow_node.IAction {
	response := make(chan flow_node.IAction, 1)

	msg := nextActionMessage{
		response: response,
	}

	variables := node.Locator.CloneVariables()
	headers := map[string]any{}
	dataSets := map[string]any{}
	if extension := node.element.ExtensionElementsField; extension != nil {
		if properties := extension.PropertiesField; properties != nil {
			fields := properties.ItemFields
			for _, field := range fields {
				value := field.ValueFor()
				if len(strings.TrimSpace(field.Value)) == 0 {
					value = variables[field.Key]
				}
				dataSets[field.Key] = value
			}
		}
		if header := extension.TaskHeaderField; header != nil {
			fields := header.ItemFields
			for _, field := range fields {
				value := field.ValueFor()
				headers[field.Key] = value
			}
		}
	}
	msg.DataObjects = node.Locator.CloneItems("$")
	msg.Headers = headers
	msg.Properties = dataSets

	node.runnerChannel <- msg
	return response
}

func (node *ScriptTask) Element() schema.FlowNodeInterface {
	return node.element
}

func (node *ScriptTask) Cancel() <-chan bool {
	response := make(chan bool)
	node.runnerChannel <- cancelMessage{response: response}
	return response
}
