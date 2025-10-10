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

func TestStartEvent(t *testing.T) {
	var testDoc schema.Definitions

	LoadTestFile("testdata/start.bpmn", &testDoc)

	engine := bpmn.NewEngine()
	proc, err := engine.NewProcess(&testDoc)
	if err == nil {
		traces := proc.Tracer().Subscribe()
		err := proc.StartAll(context.TODO())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			var trace tracing.ITrace

			select {
			case trace = <-traces:
				trace = tracing.Unwrap(trace)
			default:
				continue
			}
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
				if id, present := trace.Source.Id(); present {
					if *id == "start" {
						// success!
						t.Logf("do start event!!")
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
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
