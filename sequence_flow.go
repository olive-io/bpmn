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
	"fmt"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/errors"
)

type SequenceFlow struct {
	*schema.SequenceFlow
	process schema.Element
}

func MakeSequenceFlow(sf *schema.SequenceFlow, process schema.Element) SequenceFlow {
	return SequenceFlow{
		SequenceFlow: sf,
		process:      process,
	}
}

func NewSequenceFlow(sf *schema.SequenceFlow, process schema.Element) *SequenceFlow {
	seqFlow := MakeSequenceFlow(sf, process)
	return &seqFlow
}

func (sf *SequenceFlow) resolveId(id *string) (result schema.FlowNodeInterface, err error) {
	process := sf.process

	predicate := schema.ExactId(*id).And(schema.ElementInterface((*schema.FlowNodeInterface)(nil)))
	if flowNode, found := process.FindBy(predicate); found {
		result = flowNode.(schema.FlowNodeInterface)
	} else {
		err = errors.NotFoundError{Expected: fmt.Sprintf("flow node with ID %s", *id)}
	}
	return
}

func (sf *SequenceFlow) Source() (schema.FlowNodeInterface, error) {
	return sf.resolveId(sf.SequenceFlow.SourceRef())
}

func (sf *SequenceFlow) Target() (schema.FlowNodeInterface, error) {
	return sf.resolveId(sf.SequenceFlow.TargetRef())
}

func (sf *SequenceFlow) TargetIndex() (index int, err error) {
	var target schema.FlowNodeInterface
	target, err = sf.Target()
	if err != nil {
		return
	}
	// ownId is present since Target() already checked for this
	ownId, _ := sf.SequenceFlow.Id()
	incomings := target.Incomings()
	for i := range *incomings {
		if string((*incomings)[i]) == *ownId {
			index = i
			return
		}
	}
	err = errors.NotFoundError{Expected: fmt.Sprintf("matching incoming for %s", *ownId)}
	return
}
