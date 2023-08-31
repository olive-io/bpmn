package sequence_flow

import (
	"fmt"

	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/schema"
)

type SequenceFlow struct {
	*schema.SequenceFlow
	definitions *schema.Definitions
}

func Make(sequenceFlow *schema.SequenceFlow, definitions *schema.Definitions) SequenceFlow {
	return SequenceFlow{
		SequenceFlow: sequenceFlow,
		definitions:  definitions,
	}
}

func New(sequenceFlow *schema.SequenceFlow, definitions *schema.Definitions) *SequenceFlow {
	seqFlow := Make(sequenceFlow, definitions)
	return &seqFlow
}

func (sequenceFlow *SequenceFlow) resolveId(id *string) (result schema.FlowNodeInterface, err error) {
	ownId, present := sequenceFlow.SequenceFlow.Id()
	if !present {
		err = errors.InvalidStateError{
			Expected: "SequenceFlow to have an FlowNodeId",
			Actual:   "FlowNodeId is not present",
		}
		return
	}
	var process *schema.Process
	for i := range *sequenceFlow.definitions.Processes() {
		proc := &(*sequenceFlow.definitions.Processes())[i]
		sequenceFlows := proc.SequenceFlows()
		for j := range *sequenceFlows {
			if idPtr, present := (*sequenceFlows)[j].Id(); present {
				if *idPtr == *ownId {
					process = proc
				}
			}
		}
	}
	if process == nil {
		err = errors.NotFoundError{
			Expected: fmt.Sprintf("sequence flow with ID %s", *ownId),
		}
		return
	}

	predicate := schema.ExactId(*id).And(schema.ElementInterface((*schema.FlowNodeInterface)(nil)))
	if flowNode, found := process.FindBy(predicate); found {
		result = flowNode.(schema.FlowNodeInterface)
	} else {
		err = errors.NotFoundError{Expected: fmt.Sprintf("flow node with ID %s", *id)}
	}
	return
}

func (sequenceFlow *SequenceFlow) Source() (schema.FlowNodeInterface, error) {
	return sequenceFlow.resolveId(sequenceFlow.SequenceFlow.SourceRef())
}

func (sequenceFlow *SequenceFlow) Target() (schema.FlowNodeInterface, error) {
	return sequenceFlow.resolveId(sequenceFlow.SequenceFlow.TargetRef())
}

func (sequenceFlow *SequenceFlow) TargetIndex() (index int, err error) {
	var target schema.FlowNodeInterface
	target, err = sequenceFlow.Target()
	if err != nil {
		return
	}
	// ownId is present since Target() already checked for this
	ownId, _ := sequenceFlow.SequenceFlow.Id()
	incomings := target.Incomings()
	for i := range *incomings {
		if (*incomings)[i] == *ownId {
			index = i
			return
		}
	}
	err = errors.NotFoundError{Expected: fmt.Sprintf("matching incoming for %s", *ownId)}
	return
}
