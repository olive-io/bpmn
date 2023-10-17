// Copyright 2023 Lack (xingyys@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package receive_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity/receive"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
)

//go:embed testdata
var testdata embed.FS

func LoadTestFile(filename string, definitions any) {
	var err error
	src, err := testdata.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't read file %s: %v", filename, err)
	}
	err = xml.Unmarshal(src, definitions)
	if err != nil {
		log.Fatalf("XML unmarshalling error in %s: %v", filename, err)
	}
}

var testTask schema.Definitions

func init() {
	LoadTestFile("testdata/receive_task.bpmn", &testTask)
}

func TestReceiveTask(t *testing.T) {
	processElement := (*testTask.Processes())[0]
	proc := process.New(&processElement, &testTask)
	option := instance.WithVariables(map[string]any{
		"c": map[string]string{"name": "cc"},
	})
	if instance, err := proc.Instantiate(option); err == nil {
		traces := instance.Tracer.Subscribe()
		err := instance.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case flow.Trace:
				if id, present := trace.Source.Id(); present {
					if *id == "task" {
						// success!
						//break loop
					}
				}
			case *receive.ActiveTrace:
				trace.Do(receive.WithProperties(map[string]any{"a": "b"}))
				t.Logf("%#v", trace)
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			case flow.CeaseFlowTrace:
				break loop
			default:
				t.Logf("%#v", trace)
			}
		}
		instance.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
