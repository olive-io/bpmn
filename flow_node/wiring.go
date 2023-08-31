package flow_node

import (
	"fmt"
	"sync"

	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/id"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/sequence_flow"
	"github.com/olive-io/bpmn/tracing"
)

// Wiring holds all necessary "wiring" for functioning of
// flow nodes: definitions, process, sequence flow, event management,
// tracer, flow node mapping and a flow wait group
type Wiring struct {
	ProcessInstanceId              id.Id
	FlowNodeId                     schema.Id
	Definitions                    *schema.Definitions
	Incoming                       []sequence_flow.SequenceFlow
	Outgoing                       []sequence_flow.SequenceFlow
	EventIngress                   event.IConsumer
	EventEgress                    event.ISource
	Tracer                         tracing.ITracer
	Process                        *schema.Process
	FlowNodeMapping                *FlowNodeMapping
	FlowWaitGroup                  *sync.WaitGroup
	EventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder
}

func sequenceFlows(process *schema.Process,
	definitions *schema.Definitions,
	flows *[]schema.QName) (result []sequence_flow.SequenceFlow, err error) {
	result = make([]sequence_flow.SequenceFlow, len(*flows))
	for i := range result {
		identifier := (*flows)[i]
		exactId := schema.ExactId(identifier)
		if element, found := process.FindBy(func(e schema.Element) bool {
			_, ok := e.(*schema.SequenceFlow)
			return ok && exactId(e)
		}); found {
			result[i] = sequence_flow.Make(element.(*schema.SequenceFlow), definitions)
		} else {
			err = errors.NotFoundError{Expected: identifier}
			return
		}
	}
	return
}

func NewWiring(
	processInstanceId id.Id,
	process *schema.Process,
	definitions *schema.Definitions,
	flowNode *schema.FlowNode,
	eventIngress event.IConsumer,
	eventEgress event.ISource,
	tracer tracing.ITracer,
	flowNodeMapping *FlowNodeMapping,
	flowWaitGroup *sync.WaitGroup,
	eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder,
) (node *Wiring, err error) {
	incoming, err := sequenceFlows(process, definitions, flowNode.Incomings())
	if err != nil {
		return
	}
	outgoing, err := sequenceFlows(process, definitions, flowNode.Outgoings())
	if err != nil {
		return
	}
	var ownId string
	if ownIdPtr, present := flowNode.Id(); !present {
		err = errors.NotFoundError{
			Expected: fmt.Sprintf("flow node %#v to have an ID", flowNode),
		}
		return
	} else {
		ownId = *ownIdPtr
	}
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
	}
	return
}

// CloneFor copies receiver, overriding FlowNodeId, Incoming, Outgoing for a given flowNode
func (wiring *Wiring) CloneFor(flowNode *schema.FlowNode) (result *Wiring, err error) {
	incoming, err := sequenceFlows(wiring.Process, wiring.Definitions, flowNode.Incomings())
	if err != nil {
		return
	}
	outgoing, err := sequenceFlows(wiring.Process, wiring.Definitions, flowNode.Outgoings())
	if err != nil {
		return
	}
	var ownId string
	if ownIdPtr, present := flowNode.Id(); !present {
		err = errors.NotFoundError{
			Expected: fmt.Sprintf("flow node %#v to have an ID", flowNode),
		}
		return
	} else {
		ownId = *ownIdPtr
	}
	result = &Wiring{
		ProcessInstanceId:              wiring.ProcessInstanceId,
		FlowNodeId:                     ownId,
		Definitions:                    wiring.Definitions,
		Incoming:                       incoming,
		Outgoing:                       outgoing,
		EventIngress:                   wiring.EventIngress,
		EventEgress:                    wiring.EventEgress,
		Tracer:                         wiring.Tracer,
		Process:                        wiring.Process,
		FlowNodeMapping:                wiring.FlowNodeMapping,
		FlowWaitGroup:                  wiring.FlowWaitGroup,
		EventDefinitionInstanceBuilder: wiring.EventDefinitionInstanceBuilder,
	}
	return
}
