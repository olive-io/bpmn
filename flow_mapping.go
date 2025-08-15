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

func (mapping *FlowNodeMapping) RegisterElementToFlowNode(element schema.FlowNodeInterface, flowNode IFlowNode) error {
	id, present := element.Id()
	if !present {
		return errors.RequirementExpectationError{
			Expected: "All flow nodes must have an ID",
			Actual:   element,
		}
	}
	mapping.mapping[*id] = flowNode
	return nil
}

func (mapping *FlowNodeMapping) Finalize() {
	mapping.lock.Unlock()
}

func (mapping *FlowNodeMapping) ResolveElementToFlowNode(element schema.FlowNodeInterface) (flowNode IFlowNode, found bool) {
	mapping.lock.RLock()
	if id, present := element.Id(); present {
		flowNode, found = mapping.mapping[*id]
	}
	mapping.lock.RUnlock()
	return
}
