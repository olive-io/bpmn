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

package logic

import (
	"testing"

	"github.com/olive-io/bpmn/schema"
	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/event"
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
