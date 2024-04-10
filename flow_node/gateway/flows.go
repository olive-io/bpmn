/*
   Copyright 2023 The bpmn Authors

   This library is free software; you can redistribute it and/or
   modify it under the terms of the GNU Lesser General Public
   License as published by the Free Software Foundation; either
   version 2.1 of the License, or (at your option) any later version.

   This library is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
   Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public
   License along with this library;
*/

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
