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

package clock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextKey_String(t *testing.T) {
	t.Run("standard key", func(t *testing.T) {
		key := contextKey("test-key")
		result := key.String()
		assert.Equal(t, "clock package context key test-key", result)
	})

	t.Run("empty key", func(t *testing.T) {
		key := contextKey("")
		result := key.String()
		assert.Equal(t, "clock package context key ", result)
	})

	t.Run("special characters in key", func(t *testing.T) {
		key := contextKey("test-key-123!@#")
		result := key.String()
		assert.Equal(t, "clock package context key test-key-123!@#", result)
	})
}

func TestToContext_FromContext(t *testing.T) {
	ctx := context.Background()
	mockClock := NewMock()

	t.Run("save and retrieve clock from context", func(t *testing.T) {
		newCtx := ToContext(ctx, mockClock)
		retrievedClock, err := FromContext(newCtx)
		
		require.NoError(t, err)
		assert.Equal(t, mockClock, retrievedClock)
	})

	t.Run("retrieve clock from empty context creates host clock", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		
		clock, err := FromContext(ctx)
		require.NoError(t, err)
		assert.NotNil(t, clock)
		assert.IsType(t, &host{}, clock)
	})

	t.Run("context value is properly typed", func(t *testing.T) {
		newCtx := ToContext(ctx, mockClock)
		val := newCtx.Value(contextKey("clock"))
		
		assert.NotNil(t, val)
		assert.Equal(t, mockClock, val.(IClock))
	})

	t.Run("nested context preserves clock", func(t *testing.T) {
		clockCtx := ToContext(ctx, mockClock)
		childCtx, cancel := context.WithCancel(clockCtx)
		defer cancel()
		
		retrievedClock, err := FromContext(childCtx)
		require.NoError(t, err)
		assert.Equal(t, mockClock, retrievedClock)
	})

	t.Run("overwrite existing clock in context", func(t *testing.T) {
		firstClock := NewMock()
		secondClock := NewMockAt(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
		
		firstCtx := ToContext(ctx, firstClock)
		secondCtx := ToContext(firstCtx, secondClock)
		
		retrievedClock, err := FromContext(secondCtx)
		require.NoError(t, err)
		assert.Equal(t, secondClock, retrievedClock)
		assert.NotEqual(t, firstClock, retrievedClock)
	})
}

func TestIClock_Interface(t *testing.T) {
	t.Run("Mock implements IClock", func(t *testing.T) {
		mock := NewMock()
		assert.Implements(t, (*IClock)(nil), mock)
	})

	t.Run("host implements IClock", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		
		hostClock, err := Host(ctx)
		require.NoError(t, err)
		assert.Implements(t, (*IClock)(nil), hostClock)
	})
}

func TestFromContext_EdgeCases(t *testing.T) {
	t.Run("context with non-clock value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKey("clock"), "not-a-clock")
		
		assert.Panics(t, func() {
			_, _ = FromContext(ctx)
		})
	})

	t.Run("cancelled context still returns host clock", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		
		clock, err := FromContext(ctx)
		require.NoError(t, err)
		assert.NotNil(t, clock)
	})
}

func TestHost_Methods(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	hostClock, err := Host(ctx)
	require.NoError(t, err)
	
	t.Run("Now returns current time", func(t *testing.T) {
		before := time.Now()
		now := hostClock.Now()
		after := time.Now()
		
		assert.True(t, now.Equal(before) || now.After(before))
		assert.True(t, now.Equal(after) || now.Before(after))
	})
	
	t.Run("After with zero duration", func(t *testing.T) {
		ch := hostClock.After(0)
		assert.NotNil(t, ch)
		
		select {
		case <-ch:
			// Should receive immediately
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Should have received immediately")
		}
	})
	
	t.Run("Until with past time", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour)
		ch := hostClock.Until(pastTime)
		assert.NotNil(t, ch)
		
		select {
		case <-ch:
			// Should receive immediately for past time
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Should have received immediately for past time")
		}
	})
	
	t.Run("Changes returns non-nil channel", func(t *testing.T) {
		ch := hostClock.Changes()
		assert.NotNil(t, ch)
	})
}