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

package bpmn

import (
	"github.com/olive-io/bpmn/schema"
)

type IOutgoing interface {
	NextAction(flow Flow) chan IAction
}

type IFlowNode interface {
	IOutgoing
	Element() schema.FlowNodeInterface
}

type imessage interface {
	message()
}

type nextActionMessage struct {
	response chan IAction
	flow     Flow
}

func (m nextActionMessage) message() {}

type cancelMessage struct {
	response chan bool
}

func (m cancelMessage) message() {}

func AllSequenceFlows(sequenceFlows *[]SequenceFlow, exclusion ...func(*SequenceFlow) bool) (result []*SequenceFlow) {
	result = make([]*SequenceFlow, 0)
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
