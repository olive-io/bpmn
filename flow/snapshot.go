package flow

import (
	"github.com/olive-io/bpmn/id"
	"github.com/olive-io/bpmn/sequence_flow"
)

type Snapshot struct {
	flowId       id.Id
	sequenceFlow *sequence_flow.SequenceFlow
}

func (s *Snapshot) Id() id.Id {
	return s.flowId
}

func (s *Snapshot) SequenceFlow() *sequence_flow.SequenceFlow {
	return s.sequenceFlow
}
