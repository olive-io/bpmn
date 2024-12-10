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

type CallActivityKey struct{}

type nextCallActionMessage struct {
	headers  map[string]any
	response chan IAction
}

func (m nextCallActionMessage) message() {}

type CallActivity struct {
	*Wiring
	ctx           context.Context
	cancel        context.CancelFunc
	element       *schema.CallActivity
	runnerChannel chan imessage
}

func NewCallActivity(ctx context.Context, task *schema.CallActivity) Constructor {
	return func(wiring *Wiring) (node Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		taskNode := &CallActivity{
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

func (task *CallActivity) runner(ctx context.Context) {
	for {
		select {
		case msg := <-task.runnerChannel:
			switch m := msg.(type) {
			case cancelMessage:
				task.cancel()
				m.response <- true
			case nextCallActionMessage:
				go func() {
					aResponse := &FlowActionResponse{
						Variables: map[string]data.IItem{},
					}
					action := FlowAction{
						Response:      aResponse,
						SequenceFlows: AllSequenceFlows(&task.Outgoing),
					}

					extensionElement := task.element.ExtensionElementsField
					if extensionElement == nil || extensionElement.CalledElement == nil {
						m.response <- action
						return
					}

					headers := m.headers
					timeout := fetchTaskTimeout(headers)

					at := NewTaskTraceBuilder().
						Context(task.ctx).
						Value(CallActivityKey{}, extensionElement.CalledElement).
						Timeout(timeout).
						Activity(task).
						Headers(headers).
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
						aResponse.DataObjects = ApplyTaskDataOutput(task.element, out.DataObjects)
						aResponse.Variables = ApplyTaskResult(task.element, out.Results)
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

func (task *CallActivity) NextAction(T) chan IAction {
	response := make(chan IAction, 1)

	headers, _, _ := FetchTaskDataInput(task.Locator, task.element)
	msg := nextCallActionMessage{
		headers:  headers,
		response: response,
	}

	task.runnerChannel <- msg
	return response
}

func (task *CallActivity) Element() schema.FlowNodeInterface {
	return task.element
}

func (task *CallActivity) Type() Type {
	return CallType
}

func (task *CallActivity) Cancel() <-chan bool {
	response := make(chan bool)
	task.runnerChannel <- cancelMessage{response: response}
	return response
}
