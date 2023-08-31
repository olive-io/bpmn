package process

import (
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
)

// Trace wraps any trace within a given process
type Trace struct {
	Process *schema.Process
	Trace   tracing.ITrace
}

func (t Trace) Unwrap() tracing.ITrace {
	return t.Trace
}

func (t Trace) TraceInterface() {}
