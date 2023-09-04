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
	"sync"

	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/schema"
)

type FlowNodeMapping struct {
	mapping map[string]IFlowNode
	lock    sync.RWMutex
}

func NewLockedFlowNodeMapping() *FlowNodeMapping {
	mapping := &FlowNodeMapping{
		mapping: make(map[string]IFlowNode),
		lock:    sync.RWMutex{},
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
