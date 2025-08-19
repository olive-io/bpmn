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
	"fmt"
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/id"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

// wiring holds all necessary "wiring" for functioning of
// flow nodes: definitions, process, sequence flow, event management,
// tracer, flow node mapping and a flow wait group
type wiring struct {
	processInstanceId              id.Id
	flowNodeId                     schema.Id
	definitions                    *schema.Definitions
	incoming                       []SequenceFlow
	outgoing                       []SequenceFlow
	eventIngress                   event.IConsumer
	eventEgress                    event.ISource
	tracer                         tracing.ITracer
	process                        schema.Element
	flowNodeMapping                *FlowNodeMapping
	flowWaitGroup                  *sync.WaitGroup
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder
	locator                        data.IFlowDataLocator
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

func newWiring(
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
) (node *wiring, err error) {

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
	node = &wiring{
		processInstanceId:              processInstanceId,
		flowNodeId:                     ownId,
		definitions:                    definitions,
		incoming:                       incoming,
		outgoing:                       outgoing,
		eventIngress:                   eventIngress,
		eventEgress:                    eventEgress,
		tracer:                         tracer,
		process:                        process,
		flowNodeMapping:                flowNodeMapping,
		flowWaitGroup:                  flowWaitGroup,
		eventDefinitionInstanceBuilder: eventDefinitionInstanceBuilder,
		locator:                        locator,
	}
	return
}

// CloneFor copies receiver, overriding FlowNodeId, Incoming, Outgoing for a given flowNode
func (wr *wiring) CloneFor(node *schema.FlowNode) (result *wiring, err error) {
	incoming, err := sequenceFlows(wr.process, node.Incomings())
	if err != nil {
		return
	}
	outgoing, err := sequenceFlows(wr.process, node.Outgoings())
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
	result = &wiring{
		processInstanceId:              wr.processInstanceId,
		flowNodeId:                     ownId,
		definitions:                    wr.definitions,
		incoming:                       incoming,
		outgoing:                       outgoing,
		eventIngress:                   wr.eventIngress,
		eventEgress:                    wr.eventEgress,
		tracer:                         wr.tracer,
		process:                        wr.process,
		flowNodeMapping:                wr.flowNodeMapping,
		flowWaitGroup:                  wr.flowWaitGroup,
		eventDefinitionInstanceBuilder: wr.eventDefinitionInstanceBuilder,
	}
	return
}
