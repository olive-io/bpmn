/*
Copyright 2026 The bpmn Authors

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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFallbackGenerator(t *testing.T) {
	generator := NewFallbackGenerator()
	require.NotNil(t, generator)

	id1 := generator.New()
	id2 := generator.New()

	assert.NotEqual(t, id1.String(), id2.String())
	assert.NotEmpty(t, id1.Bytes())
	assert.Contains(t, id1.String(), "fallback-")
}

func TestFallbackGeneratorSnapshot(t *testing.T) {
	generator := NewFallbackGenerator()
	_ = generator.New()
	_ = generator.New()

	snapshot, err := generator.Snapshot()
	require.NoError(t, err)
	require.NotEmpty(t, snapshot)

	decoded := make(map[string]any)
	err = json.Unmarshal(snapshot, &decoded)
	require.NoError(t, err)

	assert.Contains(t, decoded, "prefix")
	assert.Contains(t, decoded, "counter")
}
