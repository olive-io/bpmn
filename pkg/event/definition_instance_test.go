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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/olive-io/bpmn/schema"
)

func TestWrapEventDefinition(t *testing.T) {
	definition := &mockEventDefinition{eventType: "test"}

	instance := WrapEventDefinition(definition)

	assert.NotNil(t, instance)
	assert.Equal(t, definition, instance.EventDefinition())
	assert.IsType(t, &wrappedDefinitionInstance{}, instance)
}

func TestWrappedDefinitionInstance_EventDefinition(t *testing.T) {
	definition := &mockEventDefinition{eventType: "test"}
	instance := &wrappedDefinitionInstance{definition: definition}

	result := instance.EventDefinition()

	assert.Equal(t, definition, result)
}

func TestWrappingDefinitionInstanceBuilder_NewEventDefinitionInstance(t *testing.T) {
	builder := WrappingDefinitionInstanceBuilder
	definition := &mockEventDefinition{eventType: "test"}

	instance, err := builder.NewEventDefinitionInstance(definition)

	require.NoError(t, err)
	assert.NotNil(t, instance)
	assert.Equal(t, definition, instance.EventDefinition())
}

type mockDefinitionInstanceBuilder struct {
	shouldReturnNil bool
	shouldError     bool
	instance        IDefinitionInstance
	err             error
}

func (m *mockDefinitionInstanceBuilder) NewEventDefinitionInstance(def schema.EventDefinitionInterface) (IDefinitionInstance, error) {
	if m.shouldError {
		return nil, m.err
	}
	if m.shouldReturnNil {
		return nil, nil
	}
	return m.instance, nil
}

func TestFallbackDefinitionInstanceBuilder_NewEventDefinitionInstance(t *testing.T) {
	definition := &mockEventDefinition{eventType: "test"}

	t.Run("first builder succeeds", func(t *testing.T) {
		expectedInstance := &wrappedDefinitionInstance{definition: definition}
		builder1 := &mockDefinitionInstanceBuilder{instance: expectedInstance}
		builder2 := &mockDefinitionInstanceBuilder{shouldReturnNil: true}

		fallback := &fallbackDefinitionInstanceBuilder{builders: []IDefinitionInstanceBuilder{builder1, builder2}}

		instance, err := fallback.NewEventDefinitionInstance(definition)

		require.NoError(t, err)
		assert.Equal(t, expectedInstance, instance)
	})

	t.Run("first builder returns nil, second succeeds", func(t *testing.T) {
		expectedInstance := &wrappedDefinitionInstance{definition: definition}
		builder1 := &mockDefinitionInstanceBuilder{shouldReturnNil: true}
		builder2 := &mockDefinitionInstanceBuilder{instance: expectedInstance}

		fallback := &fallbackDefinitionInstanceBuilder{builders: []IDefinitionInstanceBuilder{builder1, builder2}}

		instance, err := fallback.NewEventDefinitionInstance(definition)

		require.NoError(t, err)
		assert.Equal(t, expectedInstance, instance)
	})

	t.Run("first builder errors", func(t *testing.T) {
		expectedError := errors.New("builder error")
		builder1 := &mockDefinitionInstanceBuilder{shouldError: true, err: expectedError}
		builder2 := &mockDefinitionInstanceBuilder{shouldReturnNil: true}

		fallback := &fallbackDefinitionInstanceBuilder{builders: []IDefinitionInstanceBuilder{builder1, builder2}}

		instance, err := fallback.NewEventDefinitionInstance(definition)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, instance)
	})

	t.Run("all builders return nil", func(t *testing.T) {
		builder1 := &mockDefinitionInstanceBuilder{shouldReturnNil: true}
		builder2 := &mockDefinitionInstanceBuilder{shouldReturnNil: true}

		fallback := &fallbackDefinitionInstanceBuilder{builders: []IDefinitionInstanceBuilder{builder1, builder2}}

		instance, err := fallback.NewEventDefinitionInstance(definition)

		assert.NoError(t, err)
		assert.Nil(t, instance)
	})

	t.Run("empty builders list", func(t *testing.T) {
		fallback := &fallbackDefinitionInstanceBuilder{builders: []IDefinitionInstanceBuilder{}}

		instance, err := fallback.NewEventDefinitionInstance(definition)

		assert.NoError(t, err)
		assert.Nil(t, instance)
	})
}

func TestDefinitionInstanceBuildingChain(t *testing.T) {
	definition := &mockEventDefinition{eventType: "test"}

	t.Run("creates fallback builder", func(t *testing.T) {
		builder1 := &mockDefinitionInstanceBuilder{shouldReturnNil: true}
		builder2 := &mockDefinitionInstanceBuilder{shouldReturnNil: true}

		chain := DefinitionInstanceBuildingChain(builder1, builder2)

		assert.NotNil(t, chain)
		assert.IsType(t, &fallbackDefinitionInstanceBuilder{}, chain)

		// Test that it works as expected
		instance, err := chain.NewEventDefinitionInstance(definition)
		assert.NoError(t, err)
		assert.Nil(t, instance)
	})

	t.Run("chain with successful builder", func(t *testing.T) {
		expectedInstance := &wrappedDefinitionInstance{definition: definition}
		builder1 := &mockDefinitionInstanceBuilder{shouldReturnNil: true}
		builder2 := &mockDefinitionInstanceBuilder{instance: expectedInstance}

		chain := DefinitionInstanceBuildingChain(builder1, builder2)

		instance, err := chain.NewEventDefinitionInstance(definition)
		require.NoError(t, err)
		assert.Equal(t, expectedInstance, instance)
	})

	t.Run("empty chain", func(t *testing.T) {
		chain := DefinitionInstanceBuildingChain()

		instance, err := chain.NewEventDefinitionInstance(definition)
		assert.NoError(t, err)
		assert.Nil(t, instance)
	})

	t.Run("single builder chain", func(t *testing.T) {
		expectedInstance := &wrappedDefinitionInstance{definition: definition}
		builder := &mockDefinitionInstanceBuilder{instance: expectedInstance}

		chain := DefinitionInstanceBuildingChain(builder)

		instance, err := chain.NewEventDefinitionInstance(definition)
		require.NoError(t, err)
		assert.Equal(t, expectedInstance, instance)
	})
}

func TestIDefinitionInstance_Interface(t *testing.T) {
	definition := &mockEventDefinition{eventType: "test"}

	t.Run("wrappedDefinitionInstance implements IDefinitionInstance", func(t *testing.T) {
		instance := &wrappedDefinitionInstance{definition: definition}
		assert.Implements(t, (*IDefinitionInstance)(nil), instance)
	})

	t.Run("mockDefinitionInstance implements IDefinitionInstance", func(t *testing.T) {
		instance := &mockDefinitionInstance{definition: definition}
		assert.Implements(t, (*IDefinitionInstance)(nil), instance)
	})
}

func TestIDefinitionInstanceBuilder_Interface(t *testing.T) {
	t.Run("WrappingDefinitionInstanceBuilder implements IDefinitionInstanceBuilder", func(t *testing.T) {
		assert.Implements(t, (*IDefinitionInstanceBuilder)(nil), WrappingDefinitionInstanceBuilder)
	})

	t.Run("fallbackDefinitionInstanceBuilder implements IDefinitionInstanceBuilder", func(t *testing.T) {
		fallback := &fallbackDefinitionInstanceBuilder{}
		assert.Implements(t, (*IDefinitionInstanceBuilder)(nil), fallback)
	})

	t.Run("mockDefinitionInstanceBuilder implements IDefinitionInstanceBuilder", func(t *testing.T) {
		mock := &mockDefinitionInstanceBuilder{}
		assert.Implements(t, (*IDefinitionInstanceBuilder)(nil), mock)
	})
}