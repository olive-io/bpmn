/*
   Copyright 2023 The bpmn Authors

   This program is offered under a commercial and under the AGPL license.
   For AGPL licensing, see below.

   AGPL licensing:
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package flow_node

import (
	"sync"

	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/errors"
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
