/*
Copyright 2023 The bpmn Authors

This library is free software; you can redistribute it and/or
modify it under the terms of the GNU Lesser General Public
License as published by the Free Software Foundation; either
version 2.1 of the License, or (at your option) any later version.

This library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public
License along with this library;
*/

package bpmn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestUserTask(t *testing.T) {
	var testTask schema.Definitions
	LoadTestFile("testdata/user_task.bpmn", &testTask)

	processElement := (*testTask.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testTask)
	option := bpmn.WithVariables(map[string]any{
		"c": map[string]string{"name": "cc"},
	})
	if instance, err := proc.Instantiate(option); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.TaskTrace:
				trace.Do()
				id, _ := trace.GetActivity().Element().Id()
				assert.NotEqual(t, *id, "")
				//t.Logf("task [%s] done", *id)
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			case bpmn.CeaseFlowTrace:
				break loop
			default:
				//t.Logf("%#v", trace)
			}
		}
		instance.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
