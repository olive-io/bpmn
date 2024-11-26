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
	"sync"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/errors"
)

type FlowNodeMapping struct {
	lock    sync.RWMutex
	mapping map[string]IFlowNode
}

func NewLockedFlowNodeMapping() *FlowNodeMapping {
	mapping := &FlowNodeMapping{
		lock:    sync.RWMutex{},
		mapping: make(map[string]IFlowNode),
	}
	mapping.lock.Lock()
	return mapping
}

func (mapping *FlowNodeMapping) RegisterElementToFlowNode(element schema.FlowNodeInterface,
	flowNode IFlowNode) (err error) {
	if id, present := element.Id(); present {
		mapping.mapping[*id] = flowNode
	} else {
		err = errors.RequirementExpectationError{
			Expected: "All flow nodes must have an ID",
			Actual:   element,
		}
	}
	return
}

func (mapping *FlowNodeMapping) Finalize() {
	mapping.lock.Unlock()
}

func (mapping *FlowNodeMapping) ResolveElementToFlowNode(
	element schema.FlowNodeInterface,
) (flowNode IFlowNode, found bool) {
	mapping.lock.RLock()
	if id, present := element.Id(); present {
		flowNode, found = mapping.mapping[*id]
	}
	mapping.lock.RUnlock()
	return
}
