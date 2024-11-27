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

type nextServiceActionMessage struct {
	Headers     map[string]any
	Properties  map[string]any
	DataObjects map[string]any
	response    chan IAction
}

func (m nextServiceActionMessage) message() {}

type ServiceTask struct {
	*Wiring
	ctx           context.Context
	cancel        context.CancelFunc
	element       *schema.ServiceTask
	runnerChannel chan imessage
}

func NewServiceTask(ctx context.Context, task *schema.ServiceTask) Constructor {
	return func(wiring *Wiring) (node Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		taskNode := &ServiceTask{
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

func (task *ServiceTask) runner(ctx context.Context) {
	for {
		select {
		case msg := <-task.runnerChannel:
			switch m := msg.(type) {
			case cancelMessage:
				task.cancel()
				m.response <- true
			case nextServiceActionMessage:
				go func() {
					aResponse := &FlowActionResponse{
						DataObjects: map[string]data.IItem{},
						Variables:   map[string]data.IItem{},
					}
					action := FlowAction{
						Response:      aResponse,
						SequenceFlows: AllSequenceFlows(&task.Outgoing),
					}

					response := make(chan DoResponse, 1)
					at := NewTaskTraceBuilder().
						Context(task.ctx).
						Activity(task).
						Headers(m.Headers).
						Properties(m.Properties).
						DataObjects(m.DataObjects).
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
						aResponse.DataObjects = ApplyTaskDataOutput(task.element, out.DataObjects)
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
			task.Tracer.Trace(CancellationFlowNodeTrace{Node: task.element})
			return
		}
	}
}

func (task *ServiceTask) NextAction(T) chan IAction {
	response := make(chan IAction, 1)

	msg := nextServiceActionMessage{
		response: response,
	}

	headers, properties, dataObjects := FetchTaskDataInput(task.Locator, task.element)
	msg.Headers = headers
	msg.Properties = properties
	msg.DataObjects = dataObjects

	task.runnerChannel <- msg
	return response
}

func (task *ServiceTask) Element() schema.FlowNodeInterface { return task.element }

func (task *ServiceTask) Type() Type { return ServiceType }

func (task *ServiceTask) Cancel() <-chan bool {
	response := make(chan bool)
	task.runnerChannel <- cancelMessage{response: response}
	return response
}