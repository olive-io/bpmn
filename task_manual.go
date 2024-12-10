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

type nextManualActionMessage struct {
	headers  map[string]any
	response chan IAction
}

func (m nextManualActionMessage) message() {}

type ManualTask struct {
	*Wiring
	ctx           context.Context
	cancel        context.CancelFunc
	element       *schema.ManualTask
	runnerChannel chan imessage
}

func NewManualTask(ctx context.Context, task *schema.ManualTask) Constructor {
	return func(wiring *Wiring) (node Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		taskNode := &ManualTask{
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

func (task *ManualTask) runner(ctx context.Context) {
	for {
		select {
		case msg := <-task.runnerChannel:
			switch m := msg.(type) {
			case cancelMessage:
				task.cancel()
				m.response <- true
			case nextManualActionMessage:
				go func() {
					aResponse := &FlowActionResponse{
						Variables: map[string]data.IItem{},
					}
					action := FlowAction{
						Response:      aResponse,
						SequenceFlows: AllSequenceFlows(&task.Outgoing),
					}

					headers := m.headers
					timeout := fetchTaskTimeout(headers)

					at := NewTaskTraceBuilder().
						Context(task.ctx).
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

func (task *ManualTask) NextAction(T) chan IAction {
	response := make(chan IAction, 1)

	headers, _, _ := FetchTaskDataInput(task.Locator, task.element)
	msg := nextManualActionMessage{
		headers:  headers,
		response: response,
	}

	task.runnerChannel <- msg
	return response
}

func (task *ManualTask) Element() schema.FlowNodeInterface {
	return task.element
}

func (task *ManualTask) Type() Type {
	return ManualType
}

func (task *ManualTask) Cancel() <-chan bool {
	response := make(chan bool)
	task.runnerChannel <- cancelMessage{response: response}
	return response
}
