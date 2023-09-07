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

package service_task_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
	_ "github.com/stretchr/testify/assert"
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
	LoadTestFile("testdata/service_task.bpmn", &testTask)
}

func TestServiceTask(t *testing.T) {
	processElement := (*testTask.Processes())[0]
	proc := process.New(&processElement, &testTask)
	option := instance.WithVariables(map[string]any{
		"c": map[string]string{"name": "cc"},
	})
	if ins, err := proc.Instantiate(option); err == nil {
		traces := ins.Tracer.Subscribe()
		err := ins.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		done := make(chan struct{}, 1)
		go func() {
			for {
				var trace tracing.ITrace
				select {
				case trace = <-traces:
				}

				trace = tracing.Unwrap(trace)
				switch trace := trace.(type) {
				case flow.Trace:
					if id, present := trace.Source.Id(); present {
						if *id == "task" {
							// success!
							//break loop
						}
					}
				case activity.TaskTrace:
					trace.Done()
					t.Logf("%#v", trace)
				case tracing.ErrorTrace:
					t.Fatalf("%#v", trace)
				case flow.CeaseFlowTrace:
					close(done)
					return
				default:
					t.Logf("%#v", trace)
				}
			}
		}()

		select {
		case <-done:
		}

		ins.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
