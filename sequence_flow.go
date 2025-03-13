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
