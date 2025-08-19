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

func TestEndEvent(t *testing.T) {
	var testDoc schema.Definitions

	LoadTestFile("testdata/start.bpmn", &testDoc)

	engine := bpmn.NewEngine()
	proc, err := engine.NewProcess(&testDoc)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	traces := proc.Tracer().Subscribe()
	err = proc.StartAll(ctx)
	if err != nil {
		t.Fatalf("failed to run the instance: %s", err)
	}
loop:
	for {
		trace := tracing.Unwrap(<-traces)
		switch trace := trace.(type) {
		case bpmn.CompletionTrace:
			if id, present := trace.Node.Id(); present {
				if *id == "end" {
					// success!
					t.Logf("do end event")
				}

			}
		case bpmn.ErrorTrace:
			t.Fatalf("%#v", trace)
		case bpmn.CeaseFlowTrace:
			break loop
		default:
			//t.Logf("%#v", trace)
		}
	}
	proc.Tracer().Unsubscribe(traces)
}
