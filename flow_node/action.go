package flow_node

import (
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/sequence_flow"
)

type Action interface {
	action()
}

type ProbeAction struct {
	SequenceFlows []*sequence_flow.SequenceFlow
	// ProbeReport is a function that needs to be called
	// wth sequence flow indices that have successful
	// condition expressions (or none)
	ProbeReport func([]int)
}

func (action ProbeAction) action() {}

type ActionTransformer func(sequenceFlowId *schema.IdRef, action Action) Action
type Terminate func(sequenceFlowId *schema.IdRef) chan bool

type FlowAction struct {
	SequenceFlows []*sequence_flow.SequenceFlow
	// Index of sequence flows that should flow without
	// conditionExpression being evaluated
	UnconditionalFlows []int
	// The actions produced by the targets should be processed by
	// this function
	ActionTransformer
	// If supplied channel sends a function that returns true, the flow action
	// is to be terminated if it wasn't already
	Terminate
}

func (action FlowAction) action() {}

type CompleteAction struct{}

func (action CompleteAction) action() {}

type NoAction struct{}

func (action NoAction) action() {}
