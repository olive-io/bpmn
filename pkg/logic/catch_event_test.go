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

package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/event"
)

func TestCatchEventSatisfier_MatchSingle(t *testing.T) {
	catchEvent := schema.DefaultCatchEvent()

	sig1 := schema.DefaultSignalEventDefinition()
	sig1name := schema.QName("sig1")
	sig1.SetSignalRef(&sig1name)

	catchEvent.SetSignalEventDefinitions([]schema.SignalEventDefinition{sig1})

	satisfier := NewCatchEventSatisfier(&catchEvent, event.WrappingDefinitionInstanceBuilder)

	var satisfied bool
	var chain int

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent("sig0"))
	assert.False(t, satisfied)
	assert.Equal(t, EventDidNotMatch, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig1name)))
	assert.True(t, satisfied)
	assert.Equal(t, 0, chain)
}

func TestCatchEventSatisfier_MatchMultiple(t *testing.T) {
	catchEvent := schema.DefaultCatchEvent()

	sig1 := schema.DefaultSignalEventDefinition()
	sig1name := schema.QName("sig1")
	sig1.SetSignalRef(&sig1name)

	sig2 := schema.DefaultSignalEventDefinition()
	sig2name := schema.QName("sig2")
	sig2.SetSignalRef(&sig2name)

	catchEvent.SetSignalEventDefinitions([]schema.SignalEventDefinition{sig1, sig2})

	satisfier := NewCatchEventSatisfier(&catchEvent, event.WrappingDefinitionInstanceBuilder)

	var satisfied bool
	var chain int

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent("sig0"))
	assert.False(t, satisfied)
	assert.Equal(t, EventDidNotMatch, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig1name)))
	assert.True(t, satisfied)
	assert.Equal(t, 0, chain)
	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig2name)))
	assert.True(t, satisfied)
	assert.Equal(t, 0, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent("sig0"))
	assert.False(t, satisfied)
	assert.Equal(t, EventDidNotMatch, chain)
}

func TestCatchEventSatisfier_MatchParallelMultiple(t *testing.T) {
	catchEvent := schema.DefaultCatchEvent()

	sig1 := schema.DefaultSignalEventDefinition()
	sig1name := schema.QName("sig1")
	sig1.SetSignalRef(&sig1name)

	sig2 := schema.DefaultSignalEventDefinition()
	sig2name := schema.QName("sig2")
	sig2.SetSignalRef(&sig2name)

	parallelMultiple := true
	catchEvent.SetParallelMultiple(&parallelMultiple)
	catchEvent.SetSignalEventDefinitions([]schema.SignalEventDefinition{sig1, sig2})

	satisfier := NewCatchEventSatisfier(&catchEvent, event.WrappingDefinitionInstanceBuilder)

	var satisfied bool
	var chain int

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent("sig0"))
	assert.False(t, satisfied)
	assert.Equal(t, EventDidNotMatch, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig1name)))
	assert.False(t, satisfied)
	assert.Equal(t, 0, chain)
	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig2name)))
	assert.True(t, satisfied)
	assert.Equal(t, 0, chain)

	// Let's try this again, in a different order
	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig2name)))
	assert.False(t, satisfied)
	assert.Equal(t, 0, chain)
	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig1name)))
	assert.True(t, satisfied)
	assert.Equal(t, 0, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent("sig0"))
	assert.False(t, satisfied)
	assert.Equal(t, EventDidNotMatch, chain)

	// Now, let's supply two series of matching events but coming in partial
	// sequences

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig1name)))
	assert.False(t, satisfied)
	assert.Equal(t, 0, chain)
	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig1name)))
	assert.False(t, satisfied)
	assert.Equal(t, 1, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig2name)))
	assert.True(t, satisfied)
	assert.Equal(t, 0, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig2name)))
	assert.True(t, satisfied)
	// the reason why chain here becomes 0 is that because that chain was satisfied
	// and removed, therefore this chain become indexed at 0
	assert.Equal(t, 0, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent("sig0"))
	assert.False(t, satisfied)
	assert.Equal(t, EventDidNotMatch, chain)

}

func TestCatchEventSatisfier_MatchParallelMultipleSingleEvent(t *testing.T) {
	catchEvent := schema.DefaultCatchEvent()

	sig1 := schema.DefaultSignalEventDefinition()
	sig1name := schema.QName("sig1")
	sig1.SetSignalRef(&sig1name)

	parallelMultiple := true
	catchEvent.SetParallelMultiple(&parallelMultiple)
	catchEvent.SetSignalEventDefinitions([]schema.SignalEventDefinition{sig1})

	satisfier := NewCatchEventSatisfier(&catchEvent, event.WrappingDefinitionInstanceBuilder)

	var satisfied bool
	var chain int

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent("sig0"))
	assert.False(t, satisfied)
	assert.Equal(t, EventDidNotMatch, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig1name)))
	assert.True(t, satisfied)
	assert.Equal(t, 0, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent(string(sig1name)))
	assert.True(t, satisfied)
	assert.Equal(t, 0, chain)

	satisfied, chain = satisfier.Satisfy(event.NewSignalEvent("sig0"))
	assert.False(t, satisfied)
	assert.Equal(t, EventDidNotMatch, chain)
}
