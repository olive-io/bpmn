package event

import (
	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/schema"
)

type IEvent interface {
	MatchesEventInstance(IDefinitionInstance) bool
}

// EndEvent Process has ended
type EndEvent struct {
	Element *schema.EndEvent
}

func MakeEndEvent(element *schema.EndEvent) EndEvent {
	return EndEvent{Element: element}
}

func (ev EndEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	// always false because there's no event definition that matches
	return false
}

// NoneEvent None event
type NoneEvent struct{}

func MakeNoneEvent() NoneEvent {
	return NoneEvent{}
}

func (ev NoneEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	// always false because there's no event definition that matches
	return false
}

// SignalEvent Signal event
type SignalEvent struct {
	signalRef string
	item      data.IItem
}

func MakeSignalEvent(signalRef string, items ...data.IItem) SignalEvent {
	return SignalEvent{signalRef: signalRef, item: data.ItemOrCollection(items)}
}

func NewSignalEvent(signalRef string, items ...data.IItem) *SignalEvent {
	event := MakeSignalEvent(signalRef, items)
	return &event
}

func (ev *SignalEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	definition, ok := instance.EventDefinition().(*schema.SignalEventDefinition)
	if !ok {
		return false
	}
	signalRef, present := definition.SignalRef()
	if !present {
		return false
	}
	return *signalRef == ev.signalRef
}

func (ev *SignalEvent) SignalRef() *string {
	return &ev.signalRef
}

// CancelEvent Cancellation event
type CancelEvent struct{}

func MakeCancelEvent() CancelEvent {
	return CancelEvent{}
}

func (ev CancelEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	_, ok := instance.EventDefinition().(*schema.CancelEventDefinition)
	return ok
}

// TerminateEvent Termination event
type TerminateEvent struct{}

func MakeTerminateEvent() TerminateEvent {
	return TerminateEvent{}
}

func (ev TerminateEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	_, ok := instance.EventDefinition().(*schema.TerminateEventDefinition)
	return ok
}

// CompensationEvent Compensation event
type CompensationEvent struct {
	activityRef string
}

func MakeCompensationEvent(activityRef string) CompensationEvent {
	return CompensationEvent{activityRef: activityRef}
}

func (ev *CompensationEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	// always false because there's no event definition that matches
	return false
}

func (ev *CompensationEvent) ActivityRef() *string {
	return &ev.activityRef
}

// MessageEvent Message event
type MessageEvent struct {
	messageRef   string
	operationRef *string
	item         data.IItem
}

func MakeMessageEvent(messageRef string, operationRef *string, items ...data.IItem) MessageEvent {
	return MessageEvent{
		messageRef:   messageRef,
		operationRef: operationRef,
		item:         data.ItemOrCollection(items),
	}
}

func NewMessageEvent(messageRef string, operationRef *string, items ...data.IItem) *MessageEvent {
	event := MakeMessageEvent(messageRef, operationRef, items...)
	return &event
}

func (ev *MessageEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	definition, ok := instance.EventDefinition().(*schema.MessageEventDefinition)
	if !ok {
		return false
	}
	messageRef, present := definition.MessageRef()
	if !present {
		return false
	}
	if *messageRef != ev.messageRef {
		return false
	}
	if ev.operationRef == nil {
		if _, present := definition.OperationRef(); present {
			return false
		}
		return true
	} else {
		operationRef, present := definition.OperationRef()
		if !present {
			return false
		}
		return *operationRef == *ev.operationRef
	}
}

func (ev *MessageEvent) MessageRef() *string {
	return &ev.messageRef
}

func (ev *MessageEvent) OperationRef() (result *string, present bool) {
	if ev.operationRef != nil {
		result = ev.operationRef
		present = true
	}
	return
}

// EscalationEvent Escalation event
type EscalationEvent struct {
	escalationRef string
	item          data.IItem
}

func MakeEscalationEvent(escalationRef string, items ...data.IItem) EscalationEvent {
	return EscalationEvent{escalationRef: escalationRef, item: data.ItemOrCollection(items)}
}

func (ev *EscalationEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	definition, ok := instance.EventDefinition().(*schema.EscalationEventDefinition)
	if !ok {
		return false
	}
	escalationRef, present := definition.EscalationRef()
	if !present {
		return false
	}
	return ev.escalationRef == *escalationRef
}

func (ev *EscalationEvent) EscalationRef() *string {
	return &ev.escalationRef
}

// LinkEvent Link event
type LinkEvent struct {
	sources []string
	target  *string
}

func MakeLinkEvent(sources []string, target *string) LinkEvent {
	return LinkEvent{
		sources: sources,
		target:  target,
	}
}

func (ev *LinkEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	definition, ok := instance.EventDefinition().(*schema.LinkEventDefinition)
	if !ok {
		return false
	}

	if ev.target == nil {
		if _, present := definition.Target(); present {
			return false
		}
	} else {
		target, present := definition.Target()
		if !present {
			return false
		}

		if *target != *ev.target {
			return false
		}

	}

	if definition.Sources() == nil {
		return false
	}

	sources := definition.Sources()

	if len(ev.sources) != len(*sources) {
		return false
	}

	for i := range ev.sources {
		if ev.sources[i] != (*sources)[i] {
			return false
		}
	}

	return true
}

func (ev *LinkEvent) Sources() *[]string {
	return &ev.sources
}

func (ev *LinkEvent) Target() (result *string, present bool) {
	if ev.target != nil {
		result = ev.target
		present = true
	}
	return
}

// ErrorEvent Error event
type ErrorEvent struct {
	errorRef string
	item     data.IItem
}

func MakeErrorEvent(errorRef string, items ...data.IItem) ErrorEvent {
	return ErrorEvent{errorRef: errorRef, item: data.ItemOrCollection(items)}
}

func (ev *ErrorEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	definition, ok := instance.EventDefinition().(*schema.ErrorEventDefinition)
	if !ok {
		return false
	}
	errorRef, present := definition.ErrorRef()
	if !present {
		return false
	}
	return *errorRef == ev.errorRef
}

func (ev *ErrorEvent) ErrorRef() *string {
	return &ev.errorRef
}

// TimerEvent represents an event that occurs when a certain timer
// is triggered.
type TimerEvent struct {
	instance IDefinitionInstance
}

func MakeTimerEvent(instance IDefinitionInstance) TimerEvent {
	return TimerEvent{instance: instance}
}

func (ev TimerEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	return instance == ev.instance
}

func (ev TimerEvent) Instance() IDefinitionInstance {
	return ev.instance
}

// ConditionalEvent represents an event that occurs when a certain timer
// is triggered.
type ConditionalEvent struct {
	instance IDefinitionInstance
}

func MakeConditionalEvent(instance IDefinitionInstance) ConditionalEvent {
	return ConditionalEvent{instance: instance}
}

func (ev *ConditionalEvent) MatchesEventInstance(instance IDefinitionInstance) bool {
	return instance == ev.instance
}

func (ev *ConditionalEvent) Instance() IDefinitionInstance {
	return ev.instance
}
