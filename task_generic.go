/*
Copyright 2023 The bpmn Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bpmn

import (
	"context"
	"sync/atomic"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
)

type nextTaskActionMessage struct {
	headers     map[string]string
	properties  map[string]any
	dataObjects map[string]any
	response    chan IAction
}

func (m nextTaskActionMessage) message() {}

type genericTask struct {
	*wiring
	element      schema.FlowNodeInterface
	activityType ActivityType
	active       atomic.Int32
	mch          chan imessage
}

func newTask(element schema.FlowNodeInterface, activityType ActivityType) constructor {
	return func(wr *wiring) (activity Activity, err error) {
		taskNode := &genericTask{
			wiring:       wr,
			element:      element,
			activityType: activityType,
			active:       atomic.Int32{},
			mch:          make(chan imessage, len(wr.incoming)*2+1),
		}
		activity = taskNode
		return
	}
}

func (task *genericTask) run(ctx context.Context) {
	for {
		select {
		case msg := <-task.mch:
			switch m := msg.(type) {
			case cancelMessage:
				task.tracer.Send(CancellationFlowNodeTrace{Node: task.element})
				if task.active.Load() > 1 {
					m.response <- false
				} else {
					m.response <- true
					task.active.Swap(0)
					return
				}
			case nextTaskActionMessage:
				go func() {
					task.active.Add(1)
					defer task.active.Add(-1)

					rsp := &FlowActionResponse{
						variables: map[string]data.IItem{},
					}

					headers := m.headers
					properties := m.properties
					dataObjects := m.dataObjects

					timeout := FetchTaskTimeout(task.element)

					at := newTaskTraceBuilder().
						Context(ctx).
						Activity(task).
						Timeout(timeout).
						Headers(headers).
						Properties(properties).
						DataObjects(dataObjects).
						Build()

					task.tracer.Send(at)
					select {
					case <-ctx.Done():
						task.tracer.Send(CancellationFlowNodeTrace{Node: task.element})
						return
					case out := <-at.out():
						rsp.err = out.Err
						rsp.dataObjects = ApplyTaskDataOutput(task.element, out.DataObjects)
						rsp.variables = ApplyTaskResult(task.element, out.Results)
						rsp.handler = out.HandlerCh
					}

					m.response <- flowAction{
						response:      rsp,
						sequenceFlows: allSequenceFlows(&task.outgoing),
					}
				}()
			default:
			}
		case <-ctx.Done():
			task.tracer.Send(CancellationFlowNodeTrace{Node: task.element})
			return
		}
	}
}

func (task *genericTask) NextAction(ctx context.Context, flow Flow) chan IAction {
	if task.active.CompareAndSwap(0, 1) {
		go task.run(ctx)
	}

	response := make(chan IAction, 1)

	headers, properties, dataObjects := FetchTaskDataInput(task.locator, task.element)
	msg := nextTaskActionMessage{
		headers:     headers,
		properties:  properties,
		dataObjects: dataObjects,
		response:    response,
	}

	task.mch <- msg
	return response
}

func (task *genericTask) Element() schema.FlowNodeInterface {
	return task.element
}

func (task *genericTask) Type() ActivityType {
	return task.activityType
}

func (task *genericTask) Cancel() <-chan bool {
	response := make(chan bool)
	task.mch <- cancelMessage{response: response}
	return response
}
