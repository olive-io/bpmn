package catch

import (
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/schema"
)

type ActiveListeningTrace struct {
	Node *schema.CatchEvent
}

func (t ActiveListeningTrace) TraceInterface() {}

// EventObservedTrace signals the fact that a particular event
// has been in fact observed by the node
type EventObservedTrace struct {
	Node  *schema.CatchEvent
	Event event.Event
}

func (t EventObservedTrace) TraceInterface() {}
