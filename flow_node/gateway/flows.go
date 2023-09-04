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

package gateway

import (
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/sequence_flow"
)

func DistributeFlows(awaitingActions []chan flow_node.IAction, sequenceFlows []*sequence_flow.SequenceFlow) {
	indices := make([]int, len(sequenceFlows))
	for i := range indices {
		indices[i] = i
	}

	for i, action := range awaitingActions {
		rangeEnd := i + 1

		// If this is a last channel awaiting action
		if rangeEnd == len(awaitingActions) {
			// give it the remainder of sequence flows
			rangeEnd = len(sequenceFlows)
		}

		if rangeEnd <= len(sequenceFlows) {
			action <- flow_node.FlowAction{
				SequenceFlows:      sequenceFlows[i:rangeEnd],
				UnconditionalFlows: indices[0 : rangeEnd-i],
			}
		} else {
			// signal completion to flows that aren't
			// getting any flows
			action <- flow_node.CompleteAction{}
		}
	}
}
