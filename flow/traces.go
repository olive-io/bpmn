package flow

import (
	"github.com/olive-io/bpmn/id"
	"github.com/olive-io/bpmn/schema"
)

type NewFlowTrace struct {
	FlowId id.Id
}

func (t NewFlowTrace) TraceInterface() {}

type Trace struct {
	Source schema.FlowNodeInterface
	Flows  []Snapshot
}

func (t Trace) TraceInterface() {}

type TerminationTrace struct {
	FlowId id.Id
	Source schema.FlowNodeInterface
}

func (t TerminationTrace) TraceInterface() {}

type CancellationTrace struct {
	FlowId id.Id
}

func (t CancellationTrace) TraceInterface() {}

type CompletionTrace struct {
	Node schema.FlowNodeInterface
}

func (t CompletionTrace) TraceInterface() {}

type CeaseFlowTrace struct{}

func (t CeaseFlowTrace) TraceInterface() {}

type VisitTrace struct {
	Node schema.FlowNodeInterface
}

func (t VisitTrace) TraceInterface() {}
