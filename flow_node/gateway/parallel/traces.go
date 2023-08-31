package parallel

import (
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/schema"
)

// IncomingFlowProcessedTrace signals that a particular flow
// has been processed. If any action have been taken, it already happened
type IncomingFlowProcessedTrace struct {
	Node *schema.ParallelGateway
	Flow flow_interface.T
}

func (t IncomingFlowProcessedTrace) TraceInterface() {}
