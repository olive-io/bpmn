package flow_node

import (
	"sync"

	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/schema"
)

type FlowNodeMapping struct {
	mapping map[string]FlowNodeInterface
	lock    sync.RWMutex
}

func NewLockedFlowNodeMapping() *FlowNodeMapping {
	mapping := &FlowNodeMapping{
		mapping: make(map[string]FlowNodeInterface),
		lock:    sync.RWMutex{},
	}
	mapping.lock.Lock()
	return mapping
}

func (mapping *FlowNodeMapping) RegisterElementToFlowNode(element schema.FlowNodeInterface,
	flowNode FlowNodeInterface) (err error) {
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
) (flowNode FlowNodeInterface, found bool) {
	mapping.lock.RLock()
	if id, present := element.Id(); present {
		flowNode, found = mapping.mapping[*id]
	}
	mapping.lock.RUnlock()
	return
}
