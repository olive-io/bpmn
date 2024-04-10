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

package activity

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/event/catch"
	"github.com/olive-io/bpmn/pkg/id"
	"github.com/olive-io/bpmn/tracing"
)

type imessage interface {
	message()
}

type nextActionMessage struct {
	flow     flow_interface.T
	response chan chan flow_node.IAction
}

func (m nextActionMessage) message() {}

type Harness struct {
	*flow_node.Wiring
	element            schema.FlowNodeInterface
	runnerChannel      chan imessage
	activity           Activity
	active             int32
	cancellation       sync.Once
	eventConsumers     []event.IConsumer
	eventConsumersLock sync.RWMutex
}

func (node *Harness) ConsumeEvent(ev event.IEvent) (result event.ConsumptionResult, err error) {
	node.eventConsumersLock.RLock()
	defer node.eventConsumersLock.RUnlock()
	if atomic.LoadInt32(&node.active) == 1 {
		result, err = event.ForwardEvent(ev, &node.eventConsumers)
	}
	return
}

func (node *Harness) RegisterEventConsumer(consumer event.IConsumer) (err error) {
	node.eventConsumersLock.Lock()
	defer node.eventConsumersLock.Unlock()
	node.eventConsumers = append(node.eventConsumers, consumer)
	return
}

func (node *Harness) Activity() Activity {
	return node.activity
}

type Constructor = func(*flow_node.Wiring) (node Activity, err error)

func NewHarness(ctx context.Context,
	wiring *flow_node.Wiring,
	element *schema.FlowNode,
	idGenerator id.IGenerator,
	constructor Constructor,
) (node *Harness, err error) {
	var activity Activity
	activity, err = constructor(wiring)
	if err != nil {
		return
	}

	boundaryEvents := make([]*schema.BoundaryEvent, 0)

	switch process := wiring.Process.(type) {
	case *schema.Process:
		for i := range *process.BoundaryEvents() {
			boundaryEvent := &(*process.BoundaryEvents())[i]
			if string(*boundaryEvent.AttachedToRef()) == wiring.FlowNodeId {
				boundaryEvents = append(boundaryEvents, boundaryEvent)
			}
		}
	case *schema.SubProcess:
		for i := range *process.BoundaryEvents() {
			boundaryEvent := &(*process.BoundaryEvents())[i]
			if string(*boundaryEvent.AttachedToRef()) == wiring.FlowNodeId {
				boundaryEvents = append(boundaryEvents, boundaryEvent)
			}
		}
	}

	node = &Harness{
		Wiring:        wiring,
		element:       element,
		runnerChannel: make(chan imessage, len(wiring.Incoming)*2+1),
		activity:      activity,
	}

	err = node.EventEgress.RegisterEventConsumer(node)
	if err != nil {
		return
	}

	for i := range boundaryEvents {
		boundaryEvent := boundaryEvents[i]
		var catchEventFlowNode *flow_node.Wiring
		catchEventFlowNode, err = wiring.CloneFor(&boundaryEvent.FlowNode)
		if err != nil {
			return
		}
		// this node becomes event egress
		catchEventFlowNode.EventEgress = node

		var catchEvent *catch.Node
		catchEvent, err = catch.New(ctx, catchEventFlowNode, &boundaryEvent.CatchEvent)
		if err != nil {
			return
		}

		var actionTransformer flow_node.ActionTransformer
		if boundaryEvent.CancelActivity() {
			actionTransformer = func(sequenceFlowId *schema.IdRef, action flow_node.IAction) flow_node.IAction {
				node.cancellation.Do(func() {
					<-node.activity.Cancel()
				})
				return action
			}
		}
		flowable := flow.New(node.Definitions, catchEvent, node.Tracer,
			node.FlowNodeMapping, node.FlowWaitGroup, idGenerator, actionTransformer, node.Locator)
		flowable.Start(ctx)
	}
	sender := node.Tracer.RegisterSender()
	go node.runner(ctx, sender)
	return
}

func (node *Harness) runner(ctx context.Context, sender tracing.ISenderHandle) {
	defer sender.Done()

	for {
		select {
		case msg := <-node.runnerChannel:
			switch m := msg.(type) {
			case nextActionMessage:
				atomic.StoreInt32(&node.active, 1)
				node.Tracer.Trace(ActiveBoundaryTrace{Start: true, Node: node.activity.Element()})
				in := node.activity.NextAction(m.flow)
				out := make(chan flow_node.IAction, 1)
				go func(ctx context.Context) {
					select {
					case out <- <-in:
						atomic.StoreInt32(&node.active, 0)
						node.Tracer.Trace(ActiveBoundaryTrace{Start: false, Node: node.activity.Element()})
					case <-ctx.Done():
						return
					}
				}(ctx)
				m.response <- out
			default:
			}
		case <-ctx.Done():
			node.Tracer.Trace(flow_node.CancellationTrace{Node: node.element})
			return
		}
	}
}

func (node *Harness) NextAction(flow flow_interface.T) chan flow_node.IAction {
	response := make(chan chan flow_node.IAction, 1)
	node.runnerChannel <- nextActionMessage{flow: flow, response: response}
	return <-response
}

func (node *Harness) Element() schema.FlowNodeInterface {
	return node.element
}
