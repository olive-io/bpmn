package instance

import (
	"github.com/olive-io/bpmn/id"
	"github.com/olive-io/bpmn/tracing"
)

// InstantiationTrace denotes instantiation of a given process
type InstantiationTrace struct {
	InstanceId id.Id
}

func (i InstantiationTrace) TraceInterface() {}

// Trace wraps any trace with process instance id
type Trace struct {
	InstanceId id.Id
	Trace      tracing.Trace
}

func (t Trace) Unwrap() tracing.Trace {
	return t.Trace
}

func (t Trace) TraceInterface() {}
