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

package bpmn

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/olive-io/bpmn/schema"
)

func TestNewProcessSetKeepsDistinctExecutableProcessPointers(t *testing.T) {
	definitions := schema.DefaultDefinitions()

	p1 := schema.NewProcessBuilder().Out()
	p1.IdField = schema.NewStringP("process-1")
	p1.IsExecutableField = schema.NewBoolP(true)

	p2 := schema.NewProcessBuilder().Out()
	p2.IdField = schema.NewStringP("process-2")
	p2.IsExecutableField = schema.NewBoolP(true)

	definitions.ProcessField = []schema.Process{*p1, *p2}

	engine := NewEngine()
	ps, err := engine.NewProcessSet(&definitions)
	require.NoError(t, err)
	require.Len(t, ps.executes, 2)

	id1, ok1 := ps.executes[0].element.Id()
	id2, ok2 := ps.executes[1].element.Id()

	require.True(t, ok1)
	require.True(t, ok2)
	assert.Equal(t, "process-1", *id1)
	assert.Equal(t, "process-2", *id2)
	assert.NotSame(t, ps.executes[0].element, ps.executes[1].element)
}
