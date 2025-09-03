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

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestConsumptionResult_Constants(t *testing.T) {
	assert.Equal(t, 0, int(Consumed))
	assert.Equal(t, 1, int(ConsumptionError))
	assert.Equal(t, 2, int(PartiallyConsumed))
	assert.Equal(t, 3, int(DropEventConsumer))
}

func TestVoidConsumer_ConsumeEvent(t *testing.T) {
	consumer := VoidConsumer{}
	event := MakeNoneEvent()

	result, err := consumer.ConsumeEvent(event)

	assert.Equal(t, Consumed, result)
	assert.NoError(t, err)
}

func TestForwardEvent(t *testing.T) {
	event := MakeNoneEvent()

	t.Run("all consumers succeed", func(t *testing.T) {
		consumer1 := &mockConsumer{result: Consumed}
		consumer2 := &mockConsumer{result: Consumed}
		consumers := []IConsumer{consumer1, consumer2}

		result, err := ForwardEvent(event, &consumers)

		assert.Equal(t, Consumed, result)
		assert.NoError(t, err)
		assert.Len(t, consumer1.consumed, 1)
		assert.Len(t, consumer2.consumed, 1)
	})

	t.Run("all consumers fail", func(t *testing.T) {
		consumer1 := &mockConsumer{result: ConsumptionError, err: errors.New("error1")}
		consumer2 := &mockConsumer{result: ConsumptionError, err: errors.New("error2")}
		consumers := []IConsumer{consumer1, consumer2}

		result, err := ForwardEvent(event, &consumers)

		assert.Equal(t, ConsumptionError, result)
		assert.Error(t, err)
		multiErr, ok := err.(*multierror.Error)
		assert.True(t, ok)
		assert.Equal(t, 2, multiErr.Len())
	})

	t.Run("some consumers fail", func(t *testing.T) {
		consumer1 := &mockConsumer{result: Consumed}
		consumer2 := &mockConsumer{result: ConsumptionError, err: errors.New("error1")}
		consumers := []IConsumer{consumer1, consumer2}

		result, err := ForwardEvent(event, &consumers)

		assert.Equal(t, PartiallyConsumed, result)
		assert.Error(t, err)
		multiErr, ok := err.(*multierror.Error)
		assert.True(t, ok)
		assert.Equal(t, 1, multiErr.Len())
	})

	t.Run("consumer returns error without ConsumptionError result", func(t *testing.T) {
		consumer1 := &mockConsumer{result: Consumed, err: errors.New("ignored error")}
		consumers := []IConsumer{consumer1}

		result, err := ForwardEvent(event, &consumers)

		assert.Equal(t, Consumed, result)
		assert.NoError(t, err)
	})

	t.Run("empty consumers list", func(t *testing.T) {
		consumers := []IConsumer{}

		result, err := ForwardEvent(event, &consumers)

		assert.Equal(t, Consumed, result)
		assert.NoError(t, err)
	})
}