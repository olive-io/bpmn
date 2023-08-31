package flow_node

import (
	"github.com/olive-io/bpmn/flow/flow_interface"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/sequence_flow"
)

type IOutgoing interface {
	NextAction(flow flow_interface.T) chan IAction
}

type IFlowNode interface {
	IOutgoing
	Element() schema.FlowNodeInterface
}

func AllSequenceFlows(
	sequenceFlows *[]sequence_flow.SequenceFlow,
	exclusion ...func(*sequence_flow.SequenceFlow) bool,
) (result []*sequence_flow.SequenceFlow) {
	result = make([]*sequence_flow.SequenceFlow, 0)
sequenceFlowsLoop:
	for i := range *sequenceFlows {
		for _, exclFun := range exclusion {
			if exclFun(&(*sequenceFlows)[i]) {
				continue sequenceFlowsLoop
			}
		}
		result = append(result, &(*sequenceFlows)[i])
	}
	return
}
