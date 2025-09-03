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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

type testGeneratorBuilder struct {
	shouldError bool
}

func (t *testGeneratorBuilder) NewIdGenerator(ctx context.Context, tracer tracing.ITracer) (IGenerator, error) {
	if t.shouldError {
		return nil, assert.AnError
	}
	return &testGenerator{}, nil
}

func (t *testGeneratorBuilder) RestoreIdGenerator(ctx context.Context, serialized []byte, tracer tracing.ITracer) (IGenerator, error) {
	if t.shouldError {
		return nil, assert.AnError
	}
	return &testGenerator{serialized: serialized}, nil
}

type testGenerator struct {
	serialized []byte
	counter    int
}

func (t *testGenerator) Snapshot() ([]byte, error) {
	return []byte("test-snapshot"), nil
}

func (t *testGenerator) New() Id {
	t.counter++
	return &testId{value: t.counter}
}

type testId struct {
	value int
}

func (t *testId) Bytes() []byte {
	return []byte{byte(t.value)}
}

func (t *testId) String() string {
	return string(rune(t.value + 64))
}

func TestIGeneratorBuilder(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()

	t.Run("successful creation", func(t *testing.T) {
		builder := &testGeneratorBuilder{shouldError: false}
		
		generator, err := builder.NewIdGenerator(ctx, tracer)
		require.NoError(t, err)
		assert.NotNil(t, generator)
		assert.Implements(t, (*IGenerator)(nil), generator)
	})

	t.Run("error during creation", func(t *testing.T) {
		builder := &testGeneratorBuilder{shouldError: true}
		
		generator, err := builder.NewIdGenerator(ctx, tracer)
		assert.Error(t, err)
		assert.Nil(t, generator)
	})

	t.Run("successful restoration", func(t *testing.T) {
		builder := &testGeneratorBuilder{shouldError: false}
		serialized := []byte("test-data")
		
		generator, err := builder.RestoreIdGenerator(ctx, serialized, tracer)
		require.NoError(t, err)
		assert.NotNil(t, generator)
		assert.Implements(t, (*IGenerator)(nil), generator)
	})

	t.Run("error during restoration", func(t *testing.T) {
		builder := &testGeneratorBuilder{shouldError: true}
		serialized := []byte("test-data")
		
		generator, err := builder.RestoreIdGenerator(ctx, serialized, tracer)
		assert.Error(t, err)
		assert.Nil(t, generator)
	})
}

func TestIGenerator(t *testing.T) {
	generator := &testGenerator{}

	t.Run("snapshot", func(t *testing.T) {
		snapshot, err := generator.Snapshot()
		require.NoError(t, err)
		assert.Equal(t, []byte("test-snapshot"), snapshot)
	})

	t.Run("new id generation", func(t *testing.T) {
		id1 := generator.New()
		assert.NotNil(t, id1)
		assert.Implements(t, (*Id)(nil), id1)

		id2 := generator.New()
		assert.NotNil(t, id2)
		assert.Implements(t, (*Id)(nil), id2)

		assert.NotEqual(t, id1.String(), id2.String())
	})
}

func TestId(t *testing.T) {
	id := &testId{value: 10}

	t.Run("bytes", func(t *testing.T) {
		bytes := id.Bytes()
		assert.Equal(t, []byte{10}, bytes)
	})

	t.Run("string", func(t *testing.T) {
		str := id.String()
		assert.Equal(t, "J", str)
	})
}

func TestInterfaceCompatibility(t *testing.T) {
	ctx := context.Background()
	tracer := newMockTracer()

	t.Run("Sno implements IGeneratorBuilder", func(t *testing.T) {
		snoBuilder := GetSno()
		assert.Implements(t, (*IGeneratorBuilder)(nil), snoBuilder)

		generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
		require.NoError(t, err)
		assert.Implements(t, (*IGenerator)(nil), generator)
	})

	t.Run("SnoGenerator implements IGenerator", func(t *testing.T) {
		snoBuilder := GetSno()
		generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
		require.NoError(t, err)

		assert.Implements(t, (*IGenerator)(nil), generator)

		snapshot, err := generator.Snapshot()
		require.NoError(t, err)
		assert.NotNil(t, snapshot)

		id := generator.New()
		assert.Implements(t, (*Id)(nil), id)
	})

	t.Run("SnoId implements Id", func(t *testing.T) {
		snoBuilder := GetSno()
		generator, err := snoBuilder.NewIdGenerator(ctx, tracer)
		require.NoError(t, err)

		id := generator.New()
		assert.Implements(t, (*Id)(nil), id)

		bytes := id.Bytes()
		assert.NotNil(t, bytes)

		str := id.String()
		assert.NotEmpty(t, str)
	})
}