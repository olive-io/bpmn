package event_based

import "github.com/olive-io/bpmn/schema"

type DeterminationMadeTrace struct {
	schema.Element
}

func (trace DeterminationMadeTrace) TraceInterface() {}
