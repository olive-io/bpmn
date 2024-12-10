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

package bpmn

import (
	"context"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
)

type nextScriptActionMessage struct {
	headers    map[string]any
	properties map[string]any
	response   chan IAction
}

func (m nextScriptActionMessage) message() {}

type ScriptTask struct {
	*Wiring
	ctx           context.Context
	cancel        context.CancelFunc
	element       *schema.ScriptTask
	runnerChannel chan imessage
}

func NewScriptTask(ctx context.Context, task *schema.ScriptTask) Constructor {
	return func(wiring *Wiring) (node Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		taskNode := &ScriptTask{
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

func (task *ScriptTask) runner(ctx context.Context) {
	for {
		select {
		case msg := <-task.runnerChannel:
			switch m := msg.(type) {
			case cancelMessage:
				task.cancel()
				m.response <- true
			case nextScriptActionMessage:
				go func() {
					aResponse := &FlowActionResponse{
						Variables: map[string]data.IItem{},
					}
					action := FlowAction{
						Response:      aResponse,
						SequenceFlows: AllSequenceFlows(&task.Outgoing),
					}

					extension := task.element.ExtensionElementsField
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

					headers := m.headers
					properties := m.properties
					timeout := fetchTaskTimeout(headers)

					at := NewTaskTraceBuilder().
						Context(task.ctx).
						Timeout(timeout).
						Activity(task).
						Headers(headers).
						Properties(properties).
						Build()

					task.Tracer.Trace(at)
					select {
					case <-ctx.Done():
						task.Tracer.Trace(CancellationFlowNodeTrace{Node: task.element})
						return
					case out := <-at.out():
						if out.Err != nil {
							aResponse.Err = out.Err
						}
						for key, value := range out.Results {
							aResponse.Variables[key] = value
						}
						aResponse.Handler = out.HandlerCh
						m.response <- action
					}

				}()
			default:
			}
		case <-ctx.Done():
			task.Tracer.Trace(CancellationFlowNodeTrace{Node: task.element})
			return
		}
	}
}

func (task *ScriptTask) NextAction(t T) chan IAction {
	response := make(chan IAction, 1)

	headers, properties, _ := FetchTaskDataInput(task.Locator, task.element)
	msg := nextScriptActionMessage{
		headers:    headers,
		properties: properties,
		response:   response,
	}

	task.runnerChannel <- msg
	return response
}

func (task *ScriptTask) Element() schema.FlowNodeInterface {
	return task.element
}

func (task *ScriptTask) Type() Type {
	return ScriptType
}

func (task *ScriptTask) Cancel() <-chan bool {
	response := make(chan bool)
	task.runnerChannel <- cancelMessage{response: response}
	return response
}
