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

type nextSendActionMessage struct {
	Headers    map[string]any
	Properties map[string]any
	response   chan IAction
}

func (m nextSendActionMessage) message() {}

type SendTask struct {
	*Wiring
	ctx           context.Context
	cancel        context.CancelFunc
	element       *schema.SendTask
	runnerChannel chan imessage
}

func NewSendTask(ctx context.Context, task *schema.SendTask) Constructor {
	return func(wiring *Wiring) (node Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		taskNode := &SendTask{
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

func (task *SendTask) runner(ctx context.Context) {
	for {
		select {
		case msg := <-task.runnerChannel:
			switch m := msg.(type) {
			case cancelMessage:
				task.cancel()
				m.response <- true
			case nextSendActionMessage:
				go func() {
					aResponse := &FlowActionResponse{
						Variables: map[string]data.IItem{},
					}
					action := FlowAction{
						Response:      aResponse,
						SequenceFlows: AllSequenceFlows(&task.Outgoing),
					}

					tctx := task.ctx
					if task.element.ExtensionElementsField != nil &&
						task.element.ExtensionElementsField.TaskDefinitionField != nil {
						st := task.element.ExtensionElementsField.TaskDefinitionField.Type
						tctx = context.WithValue(tctx, TypeKey{}, st)
					}

					response := make(chan DoResponse, 1)
					at := NewTaskTraceBuilder().
						Context(tctx).
						Activity(task).
						Headers(m.Headers).
						Properties(m.Properties).
						Response(response).
						Build()

					task.Tracer.Trace(at)
					select {
					case <-ctx.Done():
						task.Tracer.Trace(CancellationFlowNodeTrace{Node: task.element})
						return
					case out := <-response:
						if out.Err != nil {
							aResponse.Err = out.Err
						}
					}

					m.response <- action
				}()
			default:
			}
		case <-ctx.Done():
			task.Tracer.Trace(CancellationFlowNodeTrace{Node: task.element})
			return
		}
	}
}

func (task *SendTask) NextAction(t T) chan IAction {
	response := make(chan IAction, 1)

	msg := nextSendActionMessage{
		response: response,
	}

	headers, properties, _ := FetchTaskDataInput(task.Locator, task.element)
	msg.Headers = headers
	msg.Properties = properties

	task.runnerChannel <- msg
	return response
}

func (task *SendTask) Element() schema.FlowNodeInterface {
	return task.element
}

func (task *SendTask) Type() Type {
	return SendType
}

func (task *SendTask) Cancel() <-chan bool {
	response := make(chan bool)
	task.runnerChannel <- cancelMessage{response: response}
	return response
}
