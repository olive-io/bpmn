package model

import (
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/schema"
)

type EventInstantiationAttemptedTrace struct {
	Event   event.Event
	Element schema.FlowNodeInterface
}

func (e EventInstantiationAttemptedTrace) TraceInterface() {}
