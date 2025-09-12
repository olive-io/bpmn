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

package bpmn_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestTask(t *testing.T) {
	var testTask schema.Definitions

	LoadTestFile("testdata/task.bpmn", &testTask)

	engine := bpmn.NewEngine()
	ctx := context.TODO()
	if instance, err := engine.NewProcess(&testTask); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll(ctx)
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
				if id, present := trace.Source.Id(); present {
					if *id == "task" {
						// success!
						break loop
					}

				}
			case bpmn.TaskTrace:
				trace.Do()
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				//t.Logf("%#v", trace)
			}
		}
		instance.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestTaskWithBuilder(t *testing.T) {
	db := schema.NewDefinitionsBuilder()
	pb := schema.NewProcessBuilder()
	task := schema.Task{}
	task.SetName(schema.NewStringP("t1"))
	pb.AddActivity(&task)
	script := schema.ScriptTask{}
	script.SetName(schema.NewStringP("s1"))
	pb.AddActivity(&script)

	db.AddProcess(*pb.Out())

	def := db.Out()

	engine := bpmn.NewEngine()
	ctx := context.TODO()
	instance, err := engine.NewProcess(def)
	if err != nil {
		t.Fatalf("failed to run the instance: %s", err)
	}
	traces := instance.Tracer().Subscribe()
	defer instance.Tracer().Unsubscribe(traces)
	err = instance.StartAll(ctx)
	if err != nil {
		t.Fatalf("failed to run the instance: %s", err)
	}

	visited := make([]string, 0)
loop:
	for {
		trace := tracing.Unwrap(<-traces)
		switch tr := trace.(type) {
		case bpmn.TaskTrace:
			name, _ := tr.GetActivity().Element().Name()
			visited = append(visited, *name)
			tr.Do()
		case bpmn.ErrorTrace:
			t.Errorf("%#v\n", trace)
		case bpmn.CeaseFlowTrace:
			break loop
		default:
		}
	}

	assert.Equal(t, []string{"t1", "s1"}, visited)
}
