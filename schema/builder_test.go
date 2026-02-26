/*
Copyright 2025 The bpmn Authors

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

package schema

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandBytes(t *testing.T) {
	t.Run("RandBytes", func(t *testing.T) {
		out := RandBytes(2)
		assert.Equal(t, len(out), 2)
		out1 := RandBytes(2)
		assert.NotEqual(t, out, out1)
		t.Log("out:", string(out))

		out2 := RandBytes(10)
		assert.NotEqual(t, out, out2)
		t.Log("out:", string(out2))
	})
}

func TestDefinitionsWithProcess(t *testing.T) {
	db := NewDefinitionsBuilder()
	pb := NewProcessBuilder()
	task := Task{}
	pb.AddActivity(&task)
	script := ScriptTask{}
	pb.AddActivity(&script)

	db.AddProcess(*pb.Out())

	db.AutoLayout(DefaultAutoLayoutConfig())
	def := db.Out()

	assert.Equal(t, len(def.ProcessField), 1)

	data, _ := xml.MarshalIndent(def, "", " ")
	t.Log(string(data))
}

func TestDefinitionsBuilderAddProcessWithNilId(t *testing.T) {
	db := NewDefinitionsBuilder()
	p := DefaultProcess()
	p.IdField = nil

	assert.NotPanics(t, func() {
		db.AddProcess(p)
	})

	def := db.Out()
	if assert.Len(t, def.ProcessField, 1) {
		assert.NotNil(t, def.ProcessField[0].IdField)
		assert.NotEqual(t, "", *def.ProcessField[0].IdField)
	}
}

func TestDefinitionsBuilderAutoLayoutDiagram(t *testing.T) {
	db := NewDefinitionsBuilder()
	pb := NewProcessBuilder()

	task := Task{}
	pb.AddActivity(&task)
	script := ScriptTask{}
	pb.AddActivity(&script)

	db.AddProcess(*pb.Out())
	db.AutoLayout(DefaultAutoLayoutConfig())
	def := db.Out()

	if assert.NotNil(t, def.DiagramField) {
		plane := def.DiagramField.BPMNPlane()
		if assert.NotNil(t, plane) {
			assert.Len(t, plane.BPMNShapeFields, 4)
			assert.Len(t, plane.BPMNEdgeFields, 3)
			for i := range plane.BPMNShapeFields {
				shape := &plane.BPMNShapeFields[i]
				if assert.NotNil(t, shape.Bounds()) {
					assert.Greater(t, float64(shape.Bounds().Width()), 0.0)
					assert.Greater(t, float64(shape.Bounds().Height()), 0.0)
				}
			}
			for i := range plane.BPMNEdgeFields {
				edge := &plane.BPMNEdgeFields[i]
				assert.Len(t, edge.WaypointField, 2)
			}
		}
	}
}

func TestDefinitionsBuilderAddTwoProcessesCreatesCollaboration(t *testing.T) {
	db := NewDefinitionsBuilder()

	p1 := NewProcessBuilder().Out()
	p2 := NewProcessBuilder().Out()

	db.AddProcess(*p1)
	db.AddProcess(*p2)
	def := db.Out()

	if assert.Len(t, def.CollaborationField, 1) {
		assert.Len(t, def.CollaborationField[0].ParticipantField, 2)
	}
}

func TestDefinitionsBuilderLayoutSpacingOptions(t *testing.T) {
	db := NewDefinitionsBuilder()

	pb := NewProcessBuilder()
	t1 := Task{}
	t2 := Task{}
	pb.AddActivity(&t1).AddActivity(&t2)

	process := pb.Out()
	db.AddProcess(*process)
	db.AutoLayout(&AutoLayoutConfig{
		StartX:     40,
		StartY:     50,
		ColumnGap:  240,
		RowGap:     140,
		ProcessGap: 200,
	})
	def := db.Out()

	require.NotNil(t, def.DiagramField)
	plane := def.DiagramField.BPMNPlane()
	require.NotNil(t, plane)
	require.Len(t, plane.BPMNShapeFields, 4)

	index := make(map[string]BPMNShape)
	for i := range plane.BPMNShapeFields {
		id, ok := plane.BPMNShapeFields[i].BpmnElement()
		require.True(t, ok)
		index[string(*id)] = plane.BPMNShapeFields[i]
	}

	startID := *process.StartEventField[0].IdField
	firstTaskID := *process.TaskField[0].IdField
	first := index[startID]
	second := index[firstTaskID]

	require.NotNil(t, first.Bounds())
	require.NotNil(t, second.Bounds())

	assert.Equal(t, 40.0, float64(first.Bounds().X()))
	assert.Equal(t, 50.0-18.0, float64(first.Bounds().Y()))
	assert.Equal(t, 40.0+240.0, float64(second.Bounds().X()))
}
