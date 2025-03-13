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
	"fmt"
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

// Wiring holds all necessary "wiring" for functioning of
// flow nodes: definitions, process, sequence flow, event management,
// tracer, flow node mapping and a flow wait group
type Wiring struct {
	ProcessInstanceId              id.Id
	FlowNodeId                     schema.Id
	Definitions                    *schema.Definitions
	Incoming                       []SequenceFlow
	Outgoing                       []SequenceFlow
	EventIngress                   event.IConsumer
	EventEgress                    event.ISource
	Tracer                         tracing.ITracer
	Process                        schema.Element
	FlowNodeMapping                *FlowNodeMapping
	FlowWaitGroup                  *sync.WaitGroup
	EventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder
	Locator                        data.IFlowDataLocator
}

func sequenceFlows(process schema.Element, flows *[]schema.QName) (result []SequenceFlow, err error) {
	result = make([]SequenceFlow, len(*flows))
	for i := range result {
		identifier := (*flows)[i]
		exactId := schema.ExactId(string(identifier))
		if element, found := process.FindBy(func(e schema.Element) bool {
			_, ok := e.(*schema.SequenceFlow)
			return ok && exactId(e)
		}); found {
			result[i] = MakeSequenceFlow(element.(*schema.SequenceFlow), process)
		} else {
			err = errors.NotFoundError{Expected: identifier}
			return
		}
	}

	return
}

func NewWiring(
	processInstanceId id.Id,
	process schema.Element,
	definitions *schema.Definitions,
	flowNode *schema.FlowNode,
	eventIngress event.IConsumer,
	eventEgress event.ISource,
	tracer tracing.ITracer,
	flowNodeMapping *FlowNodeMapping,
	flowWaitGroup *sync.WaitGroup,
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder,
	locator data.IFlowDataLocator,
) (node *Wiring, err error) {

	incoming, err := sequenceFlows(process, flowNode.Incomings())
	if err != nil {
		return
	}
	outgoing, err := sequenceFlows(process, flowNode.Outgoings())
	if err != nil {
		return
	}
	ownIdPtr, present := flowNode.Id()
	if !present {
		err = errors.NotFoundError{
			Expected: fmt.Sprintf("flow node %#v to have an ID", flowNode),
		}
		return
	}

	ownId := *ownIdPtr
	node = &Wiring{
		ProcessInstanceId:              processInstanceId,
		FlowNodeId:                     ownId,
		Definitions:                    definitions,
		Incoming:                       incoming,
		Outgoing:                       outgoing,
		EventIngress:                   eventIngress,
		EventEgress:                    eventEgress,
		Tracer:                         tracer,
		Process:                        process,
		FlowNodeMapping:                flowNodeMapping,
		FlowWaitGroup:                  flowWaitGroup,
		EventDefinitionInstanceBuilder: eventDefinitionInstanceBuilder,
		Locator:                        locator,
	}
	return
}

// CloneFor copies receiver, overriding FlowNodeId, Incoming, Outgoing for a given flowNode
func (wr *Wiring) CloneFor(node *schema.FlowNode) (result *Wiring, err error) {
	incoming, err := sequenceFlows(wr.Process, node.Incomings())
	if err != nil {
		return
	}
	outgoing, err := sequenceFlows(wr.Process, node.Outgoings())
	if err != nil {
		return
	}

	ownIdPtr, present := node.Id()
	if !present {
		err = errors.NotFoundError{
			Expected: fmt.Sprintf("flow node %#v to have an ID", node),
		}
		return
	}

	ownId := *ownIdPtr
	result = &Wiring{
		ProcessInstanceId:              wr.ProcessInstanceId,
		FlowNodeId:                     ownId,
		Definitions:                    wr.Definitions,
		Incoming:                       incoming,
		Outgoing:                       outgoing,
		EventIngress:                   wr.EventIngress,
		EventEgress:                    wr.EventEgress,
		Tracer:                         wr.Tracer,
		Process:                        wr.Process,
		FlowNodeMapping:                wr.FlowNodeMapping,
		FlowWaitGroup:                  wr.FlowWaitGroup,
		EventDefinitionInstanceBuilder: wr.EventDefinitionInstanceBuilder,
	}
	return
}
