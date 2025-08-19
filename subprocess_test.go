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

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestSubprocess(t *testing.T) {
	var testTask schema.Definitions

	LoadTestFile("testdata/subprocess.bpmn", &testTask)
	engine := bpmn.NewEngine()
	ctx := context.Background()
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

			case bpmn.TaskTrace:
				trace.Do()
				//t.Logf("%#v", trace)
			case bpmn.CeaseFlowTrace:
				break loop
			//case flow.TerminationTrace:
			//	t.Logf("%#v", trace)
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

func TestEmbedSubprocess(t *testing.T) {
	var testTask schema.Definitions

	LoadTestFile("testdata/embed-subprocess.bpmn", &testTask)

	engine := bpmn.NewEngine()
	ctx := context.Background()
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

			case bpmn.TaskTrace:
				trace.Do()
				//t.Logf("%#v", trace)
			case bpmn.CeaseFlowTrace:
				break loop
			case bpmn.TerminationTrace:
			//	t.Logf("%#v", trace)
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
