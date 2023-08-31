package activity

import "github.com/olive-io/bpmn/schema"

type ActiveBoundaryTrace struct {
	Start bool
	Node  schema.FlowNodeInterface
}

func (b ActiveBoundaryTrace) TraceInterface() {}
