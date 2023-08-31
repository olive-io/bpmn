package process

import (
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
)

// Trace wraps any trace within a given process
type Trace struct {
	Process *schema.Process
	Trace   tracing.Trace
}

func (t Trace) Unwrap() tracing.Trace {
	return t.Trace
}

func (t Trace) TraceInterface() {}
