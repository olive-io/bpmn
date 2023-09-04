// Copyright 2023 Lack (xingyys@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
