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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFanOut(t *testing.T) {
	fanOut := NewFanOut()
	assert.NotNil(t, fanOut)
	assert.Empty(t, fanOut.eventConsumers)
}

func TestFanOut_RegisterEventConsumer(t *testing.T) {
	fanOut := NewFanOut()
	consumer1 := &mockConsumer{result: Consumed}
	consumer2 := &mockConsumer{result: Consumed}

	err := fanOut.RegisterEventConsumer(consumer1)
	assert.NoError(t, err)
	assert.Len(t, fanOut.eventConsumers, 1)

	err = fanOut.RegisterEventConsumer(consumer2)
	assert.NoError(t, err)
	assert.Len(t, fanOut.eventConsumers, 2)

	assert.Equal(t, consumer1, fanOut.eventConsumers[0])
	assert.Equal(t, consumer2, fanOut.eventConsumers[1])
}

func TestFanOut_ConsumeEvent(t *testing.T) {
	event := MakeNoneEvent()

	t.Run("with no consumers", func(t *testing.T) {
		fanOut := NewFanOut()

		result, err := fanOut.ConsumeEvent(event)

		assert.Equal(t, Consumed, result)
		assert.NoError(t, err)
	})

	t.Run("with successful consumers", func(t *testing.T) {
		fanOut := NewFanOut()
		consumer1 := &mockConsumer{result: Consumed}
		consumer2 := &mockConsumer{result: Consumed}

		fanOut.RegisterEventConsumer(consumer1)
		fanOut.RegisterEventConsumer(consumer2)

		result, err := fanOut.ConsumeEvent(event)

		assert.Equal(t, Consumed, result)
		assert.NoError(t, err)
		assert.Len(t, consumer1.consumed, 1)
		assert.Len(t, consumer2.consumed, 1)
		assert.Equal(t, event, consumer1.consumed[0])
		assert.Equal(t, event, consumer2.consumed[0])
	})

	t.Run("with failing consumers", func(t *testing.T) {
		fanOut := NewFanOut()
		consumer1 := &mockConsumer{result: ConsumptionError, err: errors.New("consumer1 error")}
		consumer2 := &mockConsumer{result: ConsumptionError, err: errors.New("consumer2 error")}

		fanOut.RegisterEventConsumer(consumer1)
		fanOut.RegisterEventConsumer(consumer2)

		result, err := fanOut.ConsumeEvent(event)

		assert.Equal(t, ConsumptionError, result)
		assert.Error(t, err)
	})

	t.Run("with mixed success/failure consumers", func(t *testing.T) {
		fanOut := NewFanOut()
		consumer1 := &mockConsumer{result: Consumed}
		consumer2 := &mockConsumer{result: ConsumptionError, err: errors.New("consumer2 error")}

		fanOut.RegisterEventConsumer(consumer1)
		fanOut.RegisterEventConsumer(consumer2)

		result, err := fanOut.ConsumeEvent(event)

		assert.Equal(t, PartiallyConsumed, result)
		assert.Error(t, err)
	})
}

func TestFanOut_ConcurrentAccess(t *testing.T) {
	fanOut := NewFanOut()
	event := MakeNoneEvent()

	const numGoroutines = 10
	const numConsumers = 50
	const numEvents = 100

	var wg sync.WaitGroup

	// Add consumers concurrently
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numConsumers; j++ {
				consumer := &mockConsumer{result: Consumed}
				err := fanOut.RegisterEventConsumer(consumer)
				assert.NoError(t, err)
			}
		}()
	}

	// Consume events concurrently
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numEvents; j++ {
				result, err := fanOut.ConsumeEvent(event)
				assert.Equal(t, Consumed, result)
				assert.NoError(t, err)
			}
		}()
	}

	wg.Wait()

	// Verify final state
	fanOut.eventConsumersLock.RLock()
	consumerCount := len(fanOut.eventConsumers)
	fanOut.eventConsumersLock.RUnlock()

	assert.Equal(t, numGoroutines*numConsumers, consumerCount)
}

func TestFanOut_ThreadSafety(t *testing.T) {
	fanOut := NewFanOut()
	event := MakeNoneEvent()
	
	done := make(chan bool, 1)
	
	// Continuously add consumers
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				consumer := &mockConsumer{result: Consumed}
				fanOut.RegisterEventConsumer(consumer)
				time.Sleep(time.Microsecond)
			}
		}
	}()
	
	// Continuously consume events
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fanOut.ConsumeEvent(event)
				time.Sleep(time.Microsecond)
			}
		}
	}()
	
	// Run for a short duration
	time.Sleep(10 * time.Millisecond)
	close(done)
	
	// Test should complete without deadlock or race conditions
}

func TestFanOut_InterfaceImplementation(t *testing.T) {
	fanOut := NewFanOut()
	
	// Verify it implements IConsumer
	var consumer IConsumer = fanOut
	assert.NotNil(t, consumer)
	
	// Verify it implements ISource
	var source ISource = fanOut
	assert.NotNil(t, source)
	
	// Test the interface methods
	mockCons := &mockConsumer{result: Consumed}
	err := source.RegisterEventConsumer(mockCons)
	assert.NoError(t, err)
	
	event := MakeNoneEvent()
	result, err := consumer.ConsumeEvent(event)
	assert.Equal(t, Consumed, result)
	assert.NoError(t, err)
	assert.Len(t, mockCons.consumed, 1)
}