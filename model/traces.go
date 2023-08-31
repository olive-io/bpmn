package model

import (
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/schema"
)

type EventInstantiationAttemptedTrace struct {
	Event   event.IEvent
	Element schema.FlowNodeInterface
}

func (e EventInstantiationAttemptedTrace) TraceInterface() {}
