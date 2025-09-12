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

package id

import (
	"context"
	"testing"
	"time"

	"github.com/muyo/sno"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSno(t *testing.T) {
	snoBuilder := GetSno()
	assert.NotNil(t, snoBuilder)
	assert.IsType(t, &Sno{}, snoBuilder)
}

func TestSno_NewIdGenerator(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()
	snoBuilder := GetSno()

	generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
	require.NoError(t, err)
	assert.NotNil(t, generator)
	assert.IsType(t, &SnoGenerator{}, generator)
}

func TestSno_RestoreIdGenerator(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()
	snoBuilder := GetSno()

	t.Run("empty bytes", func(t *testing.T) {
		generator, err := snoBuilder.RestoreIdGenerator(ctx, []byte{}, tracer)
		require.NoError(t, err)
		assert.NotNil(t, generator)
		assert.IsType(t, &SnoGenerator{}, generator)
	})

	t.Run("with valid snapshot", func(t *testing.T) {
		gen, err := snoBuilder.NewIdGenerator(ctx, tracer)
		require.NoError(t, err)

		snapshot, err := gen.Snapshot()
		require.NoError(t, err)

		restoredGen, err := snoBuilder.RestoreIdGenerator(ctx, snapshot, tracer)
		require.NoError(t, err)
		assert.NotNil(t, restoredGen)
		assert.IsType(t, &SnoGenerator{}, restoredGen)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := snoBuilder.RestoreIdGenerator(ctx, []byte("invalid json"), tracer)
		assert.Error(t, err)
	})
}

func TestSnoGenerator_New(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()
	snoBuilder := GetSno()

	generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
	require.NoError(t, err)

	id1 := generator.New()
	assert.NotNil(t, id1)
	assert.IsType(t, &SnoId{}, id1)

	id2 := generator.New()
	assert.NotNil(t, id2)
	assert.IsType(t, &SnoId{}, id2)

	assert.NotEqual(t, id1.String(), id2.String())
}

func TestSnoGenerator_Snapshot(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()
	snoBuilder := GetSno()

	generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
	require.NoError(t, err)

	snapshot, err := generator.Snapshot()
	require.NoError(t, err)
	assert.NotNil(t, snapshot)
	assert.Greater(t, len(snapshot), 0)
}

func TestSnoId_String(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()
	snoBuilder := GetSno()

	generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
	require.NoError(t, err)

	id := generator.New()
	str := id.String()
	assert.NotEmpty(t, str)
	assert.IsType(t, "", str)
}

func TestSnoId_Bytes(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()
	snoBuilder := GetSno()

	generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
	require.NoError(t, err)

	id := generator.New()
	bytes := id.Bytes()
	assert.NotNil(t, bytes)
	assert.Greater(t, len(bytes), 0)
}

func TestSnoId_Consistency(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()
	snoBuilder := GetSno()

	generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
	require.NoError(t, err)

	id := generator.New()
	str1 := id.String()
	bytes1 := id.Bytes()

	str2 := id.String()
	bytes2 := id.Bytes()

	t.Log("id: " + str2)
	assert.Equal(t, str1, str2)
	assert.Equal(t, bytes1, bytes2)
}

func TestSnoGenerator_SequenceOverflow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	tracer := newMockTracer()
	snoBuilder := GetSno()

	generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
	require.NoError(t, err)

	snoGen := generator.(*SnoGenerator)
	assert.NotNil(t, snoGen.tracer)

	notification := &sno.SequenceOverflowNotification{}
	go func() {
		time.Sleep(10 * time.Millisecond)
		select {
		case <-ctx.Done():
		}
	}()

	time.Sleep(50 * time.Millisecond)
	_ = notification
}

func TestSnoGenerator_ConcurrentGeneration(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()
	snoBuilder := GetSno()

	generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
	require.NoError(t, err)

	const numGoroutines = 10
	const idsPerGoroutine = 100
	results := make(chan string, numGoroutines*idsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < idsPerGoroutine; j++ {
				id := generator.New()
				results <- id.String()
			}
		}()
	}

	seen := make(map[string]bool)
	for i := 0; i < numGoroutines*idsPerGoroutine; i++ {
		id := <-results
		assert.False(t, seen[id], "ID should be unique: %s", id)
		seen[id] = true
	}
}
