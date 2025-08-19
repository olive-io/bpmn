/*
Copyright 2023 The bpmn Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bpmn

import (
	"context"

	"github.com/olive-io/bpmn/schema"
)

type IOutgoing interface {
	NextAction(ctx context.Context, flow Flow) chan IAction
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

func allSequenceFlows(sequenceFlows *[]SequenceFlow, exclusion ...func(*SequenceFlow) bool) (result []*SequenceFlow) {
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
