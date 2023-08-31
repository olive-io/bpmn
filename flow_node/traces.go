package flow_node

import "github.com/olive-io/bpmn/schema"

type CancellationTrace struct {
	Node schema.FlowNodeInterface
}

func (t CancellationTrace) TraceInterface() {}

type NewFlowNodeTrace struct {
	Node schema.FlowNodeInterface
}

func (t NewFlowNodeTrace) TraceInterface() {}
