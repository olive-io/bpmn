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
	headers     map[string]string
	properties  map[string]any
	dataObjects map[string]any
	response    chan IAction
}

func (m nextServiceActionMessage) message() {}

type ServiceTask struct {
	*Wiring
	ctx     context.Context
	cancel  context.CancelFunc
	element *schema.ServiceTask
	mch     chan imessage
}

func NewServiceTask(ctx context.Context, task *schema.ServiceTask) Constructor {
	return func(wiring *Wiring) (node Activity, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		taskNode := &ServiceTask{
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

func (task *ServiceTask) run(ctx context.Context) {
	for {
		select {
		case msg := <-task.mch:
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

					headers := m.headers
					properties := m.properties
					dataObjects := m.dataObjects

					timeout := fetchTaskTimeout(task.element)

					at := NewTaskTraceBuilder().
						Context(task.ctx).
						Activity(task).
						Timeout(timeout).
						Headers(headers).
						Properties(properties).
						DataObjects(dataObjects).
						Build()

					task.Tracer.Trace(at)
					select {
					case <-ctx.Done():
						task.Tracer.Trace(CancellationFlowNodeTrace{Node: task.element})
						return
					case out := <-at.out():
						aResponse.Err = out.Err
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

func (task *ServiceTask) NextAction(T) chan IAction {
	response := make(chan IAction, 1)

	headers, properties, dataObjects := FetchTaskDataInput(task.Locator, task.element)
	msg := nextServiceActionMessage{
		headers:     headers,
		properties:  properties,
		dataObjects: dataObjects,
		response:    response,
	}

	task.mch <- msg
	return response
}

func (task *ServiceTask) Element() schema.FlowNodeInterface { return task.element }

func (task *ServiceTask) Type() Type { return ServiceType }

func (task *ServiceTask) Cancel() <-chan bool {
	response := make(chan bool)
	task.mch <- cancelMessage{response: response}
	return response
}
