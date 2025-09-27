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

package event

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
)

func TestEndEvent(t *testing.T) {
	element := &schema.EndEvent{}

	t.Run("MakeEndEvent", func(t *testing.T) {
		event := MakeEndEvent(element)
		assert.Equal(t, element, event.Element)
	})

	t.Run("MatchesEventInstance always returns false", func(t *testing.T) {
		event := MakeEndEvent(element)
		instance := &mockDefinitionInstance{definition: &mockEventDefinition{}}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}

func TestNoneEvent(t *testing.T) {
	t.Run("MakeNoneEvent", func(t *testing.T) {
		event := MakeNoneEvent()
		assert.IsType(t, NoneEvent{}, event)
	})

	t.Run("MatchesEventInstance always returns false", func(t *testing.T) {
		event := MakeNoneEvent()
		instance := &mockDefinitionInstance{definition: &mockEventDefinition{}}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}

func TestSignalEvent(t *testing.T) {
	signalRef := "test-signal"
	item := schema.NewValue("test-data")

	t.Run("MakeSignalEvent", func(t *testing.T) {
		event := MakeSignalEvent(signalRef, item)
		assert.Equal(t, signalRef, event.signalRef)
		assert.Equal(t, data.ItemOrCollection(item), event.item)
	})

	t.Run("NewSignalEvent", func(t *testing.T) {
		event := NewSignalEvent(signalRef, item)
		assert.NotNil(t, event)
		assert.Equal(t, signalRef, event.signalRef)
		assert.Equal(t, data.ItemOrCollection(item), event.item)
	})

	t.Run("SignalRef", func(t *testing.T) {
		event := NewSignalEvent(signalRef, item)
		result := event.SignalRef()
		assert.Equal(t, &signalRef, result)
	})

	t.Run("MatchesEventInstance with matching signal", func(t *testing.T) {
		event := NewSignalEvent(signalRef)
		signalRefValue := schema.QName(signalRef)
		definition := &schema.SignalEventDefinition{}
		definition.SetSignalRef(&signalRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with non-matching signal", func(t *testing.T) {
		event := NewSignalEvent(signalRef)
		signalRefValue := schema.QName("different-signal")
		definition := &schema.SignalEventDefinition{}
		definition.SetSignalRef(&signalRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})

	t.Run("MatchesEventInstance with wrong definition type", func(t *testing.T) {
		event := NewSignalEvent(signalRef)
		definition := &schema.MessageEventDefinition{}
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})

	t.Run("MatchesEventInstance with no signal ref", func(t *testing.T) {
		event := NewSignalEvent(signalRef)
		definition := &schema.SignalEventDefinition{}
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}

func TestCancelEvent(t *testing.T) {
	t.Run("MakeCancelEvent", func(t *testing.T) {
		event := MakeCancelEvent()
		assert.IsType(t, CancelEvent{}, event)
	})

	t.Run("MatchesEventInstance with CancelEventDefinition", func(t *testing.T) {
		event := MakeCancelEvent()
		definition := &schema.CancelEventDefinition{}
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with wrong definition type", func(t *testing.T) {
		event := MakeCancelEvent()
		definition := &schema.SignalEventDefinition{}
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}

func TestTerminateEvent(t *testing.T) {
	t.Run("MakeTerminateEvent", func(t *testing.T) {
		event := MakeTerminateEvent()
		assert.IsType(t, TerminateEvent{}, event)
	})

	t.Run("MatchesEventInstance with TerminateEventDefinition", func(t *testing.T) {
		event := MakeTerminateEvent()
		definition := &schema.TerminateEventDefinition{}
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with wrong definition type", func(t *testing.T) {
		event := MakeTerminateEvent()
		definition := &schema.SignalEventDefinition{}
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}

func TestCompensationEvent(t *testing.T) {
	activityRef := "test-activity"

	t.Run("MakeCompensationEvent", func(t *testing.T) {
		event := MakeCompensationEvent(activityRef)
		assert.Equal(t, activityRef, event.activityRef)
	})

	t.Run("ActivityRef", func(t *testing.T) {
		event := MakeCompensationEvent(activityRef)
		result := event.ActivityRef()
		assert.Equal(t, &activityRef, result)
	})

	t.Run("MatchesEventInstance always returns false", func(t *testing.T) {
		event := MakeCompensationEvent(activityRef)
		instance := &mockDefinitionInstance{definition: &mockEventDefinition{}}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}

func TestMessageEvent(t *testing.T) {
	messageRef := "test-message"
	operationRef := "test-operation"
	item := schema.NewValue("test-data")

	t.Run("MakeMessageEvent", func(t *testing.T) {
		event := MakeMessageEvent(messageRef, &operationRef, item)
		assert.Equal(t, messageRef, event.messageRef)
		assert.Equal(t, &operationRef, event.operationRef)
		assert.Equal(t, data.ItemOrCollection(item), event.item)
	})

	t.Run("NewMessageEvent", func(t *testing.T) {
		event := NewMessageEvent(messageRef, &operationRef, item)
		assert.NotNil(t, event)
		assert.Equal(t, messageRef, event.messageRef)
		assert.Equal(t, &operationRef, event.operationRef)
		assert.Equal(t, data.ItemOrCollection(item), event.item)
	})

	t.Run("MessageRef", func(t *testing.T) {
		event := NewMessageEvent(messageRef, &operationRef)
		result := event.MessageRef()
		assert.Equal(t, &messageRef, result)
	})

	t.Run("OperationRef with operation", func(t *testing.T) {
		event := NewMessageEvent(messageRef, &operationRef)
		result, present := event.OperationRef()
		assert.True(t, present)
		assert.Equal(t, &operationRef, result)
	})

	t.Run("OperationRef without operation", func(t *testing.T) {
		event := NewMessageEvent(messageRef, nil)
		result, present := event.OperationRef()
		assert.False(t, present)
		assert.Nil(t, result)
	})

	t.Run("MatchesEventInstance with matching message and operation", func(t *testing.T) {
		event := NewMessageEvent(messageRef, &operationRef)
		messageRefValue := schema.QName(messageRef)
		operationRefValue := schema.QName(operationRef)
		definition := &schema.MessageEventDefinition{}
		definition.SetMessageRef(&messageRefValue)
		definition.SetOperationRef(&operationRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with matching message but no operation", func(t *testing.T) {
		event := NewMessageEvent(messageRef, nil)
		messageRefValue := schema.QName(messageRef)
		definition := &schema.MessageEventDefinition{}
		definition.SetMessageRef(&messageRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with non-matching message", func(t *testing.T) {
		event := NewMessageEvent(messageRef, nil)
		messageRefValue := schema.QName("different-message")
		definition := &schema.MessageEventDefinition{}
		definition.SetMessageRef(&messageRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})

	t.Run("MatchesEventInstance with non-matching operation", func(t *testing.T) {
		event := NewMessageEvent(messageRef, &operationRef)
		messageRefValue := schema.QName(messageRef)
		operationRefValue := schema.QName("different-operation")
		definition := &schema.MessageEventDefinition{}
		definition.SetMessageRef(&messageRefValue)
		definition.SetOperationRef(&operationRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}

func TestEscalationEvent(t *testing.T) {
	escalationRef := "test-escalation"
	item := schema.NewValue("test-data")

	t.Run("MakeEscalationEvent", func(t *testing.T) {
		event := MakeEscalationEvent(escalationRef, item)
		assert.Equal(t, escalationRef, event.escalationRef)
		assert.Equal(t, data.ItemOrCollection(item), event.item)
	})

	t.Run("EscalationRef", func(t *testing.T) {
		event := MakeEscalationEvent(escalationRef)
		result := event.EscalationRef()
		assert.Equal(t, &escalationRef, result)
	})

	t.Run("MatchesEventInstance with matching escalation", func(t *testing.T) {
		event := MakeEscalationEvent(escalationRef)
		escalationRefValue := schema.QName(escalationRef)
		definition := &schema.EscalationEventDefinition{}
		definition.SetEscalationRef(&escalationRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with non-matching escalation", func(t *testing.T) {
		event := MakeEscalationEvent(escalationRef)
		escalationRefValue := schema.QName("different-escalation")
		definition := &schema.EscalationEventDefinition{}
		definition.SetEscalationRef(&escalationRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}

func TestErrorEvent(t *testing.T) {
	errorRef := "test-error"
	item := schema.NewValue("test-data")

	t.Run("MakeErrorEvent", func(t *testing.T) {
		event := MakeErrorEvent(errorRef, item)
		assert.Equal(t, errorRef, event.errorRef)
		assert.Equal(t, data.ItemOrCollection(item), event.item)
	})

	t.Run("ErrorRef", func(t *testing.T) {
		event := MakeErrorEvent(errorRef)
		result := event.ErrorRef()
		assert.Equal(t, &errorRef, result)
	})

	t.Run("MatchesEventInstance with matching error", func(t *testing.T) {
		event := MakeErrorEvent(errorRef)
		errorRefValue := schema.QName(errorRef)
		definition := &schema.ErrorEventDefinition{}
		definition.SetErrorRef(&errorRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with non-matching error", func(t *testing.T) {
		event := MakeErrorEvent(errorRef)
		errorRefValue := schema.QName("different-error")
		definition := &schema.ErrorEventDefinition{}
		definition.SetErrorRef(&errorRefValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}

func TestTimerEvent(t *testing.T) {
	instance := &mockDefinitionInstance{definition: &mockEventDefinition{}}

	t.Run("MakeTimerEvent", func(t *testing.T) {
		event := MakeTimerEvent(instance)
		assert.Equal(t, instance, event.instance)
	})

	t.Run("Instance", func(t *testing.T) {
		event := MakeTimerEvent(instance)
		result := event.Instance()
		assert.Equal(t, instance, result)
	})

	t.Run("MatchesEventInstance with same instance", func(t *testing.T) {
		event := MakeTimerEvent(instance)
		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with different instance", func(t *testing.T) {
		event := MakeTimerEvent(instance)
		differentInstance := &mockDefinitionInstance{definition: &mockEventDefinition{}}
		result := event.MatchesEventInstance(differentInstance)
		assert.False(t, result)
	})
}

func TestConditionalEvent(t *testing.T) {
	instance := &mockDefinitionInstance{definition: &mockEventDefinition{}}

	t.Run("MakeConditionalEvent", func(t *testing.T) {
		event := MakeConditionalEvent(instance)
		assert.Equal(t, instance, event.instance)
	})

	t.Run("Instance", func(t *testing.T) {
		event := MakeConditionalEvent(instance)
		result := event.Instance()
		assert.Equal(t, instance, result)
	})

	t.Run("MatchesEventInstance with same instance", func(t *testing.T) {
		event := MakeConditionalEvent(instance)
		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with different instance", func(t *testing.T) {
		event := MakeConditionalEvent(instance)
		differentInstance := &mockDefinitionInstance{definition: &mockEventDefinition{}}
		result := event.MatchesEventInstance(differentInstance)
		assert.False(t, result)
	})
}

func TestLinkEvent(t *testing.T) {
	sources := []string{"source1", "source2"}
	target := "target1"

	t.Run("MakeLinkEvent", func(t *testing.T) {
		event := MakeLinkEvent(sources, &target)
		assert.Equal(t, sources, event.sources)
		assert.Equal(t, &target, event.target)
	})

	t.Run("Sources", func(t *testing.T) {
		event := MakeLinkEvent(sources, &target)
		result := event.Sources()
		assert.Equal(t, &sources, result)
	})

	t.Run("Target with target", func(t *testing.T) {
		event := MakeLinkEvent(sources, &target)
		result, present := event.Target()
		assert.True(t, present)
		assert.Equal(t, &target, result)
	})

	t.Run("Target without target", func(t *testing.T) {
		event := MakeLinkEvent(sources, nil)
		result, present := event.Target()
		assert.False(t, present)
		assert.Nil(t, result)
	})

	t.Run("MatchesEventInstance with matching link", func(t *testing.T) {
		event := MakeLinkEvent(sources, &target)
		targetValue := schema.QName(target)
		sourcesValue := []schema.QName{schema.QName("source1"), schema.QName("source2")}
		definition := &schema.LinkEventDefinition{}
		definition.SetTarget(&targetValue)
		definition.SetSources(sourcesValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.True(t, result)
	})

	t.Run("MatchesEventInstance with non-matching target", func(t *testing.T) {
		event := MakeLinkEvent(sources, &target)
		targetValue := schema.QName("different-target")
		sourcesValue := []schema.QName{schema.QName("source1"), schema.QName("source2")}
		definition := &schema.LinkEventDefinition{}
		definition.SetTarget(&targetValue)
		definition.SetSources(sourcesValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})

	t.Run("MatchesEventInstance with non-matching sources", func(t *testing.T) {
		event := MakeLinkEvent(sources, &target)
		targetValue := schema.QName(target)
		sourcesValue := []schema.QName{schema.QName("different-source")}
		definition := &schema.LinkEventDefinition{}
		definition.SetTarget(&targetValue)
		definition.SetSources(sourcesValue)
		instance := &mockDefinitionInstance{definition: definition}

		result := event.MatchesEventInstance(instance)
		assert.False(t, result)
	})
}
