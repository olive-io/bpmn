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

package service_task

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
	Headers     map[string]any
	Properties  map[string]any
	DataObjects map[string]any
	response    chan flow_node.IAction
}

func (m nextActionMessage) message() {}

type cancelMessage struct {
	response chan bool
}

func (m cancelMessage) message() {}

type ServiceTask struct {
	*flow_node.Wiring
	ctx           context.Context
	cancel        context.CancelFunc
	element       *schema.ServiceTask
	runnerChannel chan message
}

func NewServiceTask(ctx context.Context, task *schema.ServiceTask) activity.Constructor {
	return func(wiring *flow_node.Wiring) (node activity.Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		done := &atomic.Bool{}
		done.Store(false)
		taskNode := &ServiceTask{
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

func (node *ServiceTask) runner(ctx context.Context) {
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
						DataObjects:   map[string]data.IItem{},
						Variables:     map[string]data.IItem{},
						SequenceFlows: flow_node.AllSequenceFlows(&node.Outgoing),
					}

					response := make(chan callResponse, 1)
					at := &ActiveTrace{
						Context:     node.ctx,
						Activity:    node,
						Headers:     m.Headers,
						Properties:  m.Properties,
						DataObjects: m.DataObjects,
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
						for name, do := range out.dataObjects {
							action.DataObjects[name] = do
						}
						for key, value := range out.result {
							action.Variables[key] = value
						}
					}

					m.response <- action
				}()
			default:
			}
		case <-ctx.Done():
			return
		}
	}
}

func (node *ServiceTask) NextAction(flow_interface.T) chan flow_node.IAction {
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
	msg.Headers = headers
	msg.Properties = dataSets
	msg.DataObjects = node.Locator.CloneItems("$")

	node.runnerChannel <- msg
	return response
}

func (node *ServiceTask) Element() schema.FlowNodeInterface {
	return node.element
}

func (node *ServiceTask) Cancel() <-chan bool {
	response := make(chan bool)
	node.runnerChannel <- cancelMessage{response: response}
	return response
}
