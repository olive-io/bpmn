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

package script

import (
	"context"

	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
)

type imessage interface {
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

type Node struct {
	*flow_node.Wiring
	ctx           context.Context
	cancel        context.CancelFunc
	element       *schema.ScriptTask
	runnerChannel chan imessage
}

func NewScriptTask(ctx context.Context, task *schema.ScriptTask) activity.Constructor {
	return func(wiring *flow_node.Wiring) (node activity.Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		taskNode := &Node{
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

func (node *Node) runner(ctx context.Context) {
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

					response := make(chan activity.DoResponse, 1)

					extension := node.element.ExtensionElementsField
					taskDef := extension.TaskDefinitionField
					if taskDef == nil {
						taskDef = &schema.TaskDefinition{Type: "expression"}
					}

					/*
						lang := *node.Definitions.ExpressionLanguage()
						engine := expression.GetEngine(ctx, lang)

						properties := make(map[string]any)
						for key, value := range m.Properties {
							properties[key] = value
						}
						for key, value := range m.DataObjects {
							properties[key] = value
						}

						if extension.ScriptField != nil {
							if !strings.HasPrefix(extension.ScriptField.Expression, "=") {
								aResponse.Err = errors.InvalidArgumentError{
									Expected: "script expression must start with '=', (like '=a+b')",
									Actual:   "(" + extension.ScriptField.Expression + ")",
								}
								m.response <- action
								return
							}
							expr := strings.TrimPrefix(extension.ScriptField.Expression, "=")
							compiled, err := engine.CompileExpression(expr)
							if err != nil {
								aResponse.Err = errors.InvalidArgumentError{
									Expected: "must be a legal expression",
									Actual:   err.Error(),
								}
								m.response <- action
								return
							}

							result, err := engine.EvaluateExpression(compiled, properties)
							if err != nil {
								aResponse.Err = errors.TaskExecError{Id: node.FlowNodeId, Reason: err.Error()}
								m.response <- action
								return
							}

							key := extension.ScriptField.Result
							noMatchedErr := errors.InvalidArgumentError{
								Expected: fmt.Sprintf("boolean result in conditionExpression (%s)", expr),
								Actual:   result,
							}

							switch extension.ScriptField.ResultType {
							case schema.ItemTypeBoolean:
								if v, ok := result.(bool); !ok {
									aResponse.Err = noMatchedErr
								} else {
									aResponse.Variables[key] = v
								}
							case schema.ItemTypeString:
								if v, ok := result.(string); !ok {
									aResponse.Err = noMatchedErr
								} else {
									aResponse.Variables[key] = v
								}
							case schema.ItemTypeInteger:
								if v, ok := result.(int); !ok {
									aResponse.Err = noMatchedErr
								} else {
									aResponse.Variables[key] = v
								}
							case schema.ItemTypeFloat:
								if v, ok := result.(float64); !ok {
									aResponse.Err = noMatchedErr
								} else {
									aResponse.Variables[key] = v
								}
							}
						}
					*/

					at := activity.NewTraceBuilder().
						Context(node.ctx).
						Activity(node).
						Headers(m.Headers).
						Properties(m.Properties).
						Response(response).
						Build()

					node.Tracer.Trace(at)
					select {
					case <-ctx.Done():
						node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
						return
					case out := <-response:
						if out.Err != nil {
							aResponse.Err = out.Err
						}
						for key, value := range out.Properties {
							aResponse.Variables[key] = value
						}
						aResponse.Handler = out.HandlerCh
						m.response <- action
					}

				}()
			default:
			}
		case <-ctx.Done():
			node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
			return
		}
	}
}

func (node *Node) NextAction(t flow_interface.T) chan flow_node.IAction {
	response := make(chan flow_node.IAction, 1)

	msg := nextActionMessage{
		response: response,
	}

	headers, properties, dataObjects := activity.FetchTaskDataInput(node.Locator, node.element)
	msg.Headers = headers
	msg.Properties = properties
	msg.DataObjects = dataObjects

	node.runnerChannel <- msg
	return response
}

func (node *Node) Element() schema.FlowNodeInterface {
	return node.element
}

func (node *Node) Type() activity.Type {
	return activity.ScriptType
}

func (node *Node) Cancel() <-chan bool {
	response := make(chan bool)
	node.runnerChannel <- cancelMessage{response: response}
	return response
}
