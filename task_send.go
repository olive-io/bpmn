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
	headers    map[string]any
	properties map[string]any
	response   chan IAction
}

func (m nextSendActionMessage) message() {}

type SendTask struct {
	*Wiring
	ctx     context.Context
	cancel  context.CancelFunc
	element *schema.SendTask
	mch     chan imessage
}

func NewSendTask(ctx context.Context, task *schema.SendTask) Constructor {
	return func(wiring *Wiring) (node Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		taskNode := &SendTask{
			Wiring:  wiring,
			ctx:     ctx,
			cancel:  cancel,
			element: task,
			mch:     make(chan imessage, len(wiring.Incoming)*2+1),
		}
		go taskNode.run(ctx)
		node = taskNode
		return
	}
}

func (task *SendTask) run(ctx context.Context) {
	for {
		select {
		case msg := <-task.mch:
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

					headers := m.headers
					properties := m.properties
					timeout := fetchTaskTimeout(headers)

					at := NewTaskTraceBuilder().
						Context(tctx).
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
						aResponse.Err = out.Err
						aResponse.Variables = ApplyTaskResult(task.element, out.Results)
						aResponse.Handler = out.HandlerCh
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

	headers, properties, _ := FetchTaskDataInput(task.Locator, task.element)
	msg := nextSendActionMessage{
		headers:    headers,
		properties: properties,
		response:   response,
	}

	task.mch <- msg
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
	task.mch <- cancelMessage{response: response}
	return response
}
