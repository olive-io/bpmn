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

package service_test

import (
	"context"
	"embed"
	"encoding/xml"
	"fmt"
	"log"
	"testing"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/flow_node/activity/service"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/assert"
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
	options := []instance.Option{
		instance.WithVariables(map[string]any{
			"c": map[string]string{"name": "cc"},
		}),
		instance.WithDataObjects(map[string]any{
			"a": struct{}{},
		}),
	}
	if ins, err := proc.Instantiate(options...); err == nil {
		traces := ins.Tracer.Subscribe()
		err = ins.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		done := make(chan struct{}, 1)
		go func() {
			defer close(done)
			for {
				var trace tracing.ITrace
				select {
				case trace = <-traces:
				}

				trace = tracing.Unwrap(trace)
				switch trace := trace.(type) {
				case flow.Trace:
				case activity.ActiveTaskTrace:
					trace.Execute()
					t.Logf("%#v", trace)
				case tracing.ErrorTrace:
					t.Errorf("%#v", trace)
					return
				case flow.CeaseFlowTrace:
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

func TestServiceTaskWithError(t *testing.T) {
	processElement := (*testTask.Processes())[0]
	proc := process.New(&processElement, &testTask)
	var te error
	if ins, err := proc.Instantiate(); err == nil {
		traces := ins.Tracer.Subscribe()
		err := ins.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		done := make(chan struct{}, 1)
		go func() {
			defer close(done)
			for {
				var trace tracing.ITrace
				select {
				case trace = <-traces:
				}

				trace = tracing.Unwrap(trace)
				switch trace := trace.(type) {
				case flow.Trace:
				case activity.ActiveTaskTrace:
					if st, ok := trace.(*service.ActiveTrace); ok {
						st.Do(nil, nil, fmt.Errorf("text error"), nil)
					}
					t.Logf("%#v", trace)
				case tracing.ErrorTrace:
					te = trace.Error
					t.Logf("%#v", trace)
				case flow.CeaseFlowTrace:
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

	assert.Equal(t, te.(errors.TaskExecError).Reason, "text error")
}

func TestServiceTaskWithRetry(t *testing.T) {
	processElement := (*testTask.Processes())[0]
	proc := process.New(&processElement, &testTask)
	var te error
	runnum := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if ins, err := proc.Instantiate(); err == nil {
		traces := ins.Tracer.Subscribe()
		err = ins.StartAll(ctx)
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		done := make(chan struct{}, 1)
		go func() {
			defer close(done)
			for {
				var trace tracing.ITrace
				select {
				case trace = <-traces:
				}

				trace = tracing.Unwrap(trace)
				switch trace := trace.(type) {
				case flow.Trace:
				case activity.ActiveTaskTrace:
					runnum += 1
					if st, ok := trace.(*service.ActiveTrace); ok {
						retry := int32(1)
						st.Do(nil, nil, fmt.Errorf("text error"), &retry)
					}
					t.Logf("%#v", trace)
				case tracing.ErrorTrace:
					te = trace.Error
					t.Logf("%#v", trace)
				case flow.CeaseFlowTrace:
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

	assert.Equal(t, runnum, 2)
	assert.Equal(t, te.(errors.TaskExecError).Reason, "text error")
}

func TestServiceTaskWithDataInput(t *testing.T) {
	task := &schema.Definitions{}
	LoadTestFile("testdata/service_task_data.bpmn", &task)
	processElement := (*task.Processes())[0]
	options := []instance.Option{
		instance.WithDataObjects(map[string]any{"DataObject_0yhrl3s": map[string]any{"a": "ac"}}),
	}
	proc := process.New(&processElement, task)
	if ins, err := proc.Instantiate(options...); err == nil {
		traces := ins.Tracer.Subscribe()
		err := ins.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		done := make(chan struct{}, 1)
		go func() {
			defer close(done)
			for {
				var trace tracing.ITrace
				select {
				case trace = <-traces:
				}

				trace = tracing.Unwrap(trace)
				switch trace := trace.(type) {
				case flow.Trace:
				case activity.ActiveTaskTrace:
					if st, ok := trace.(*service.ActiveTrace); ok {
						assert.Equal(t, st.DataObjects["in"], map[string]any{"a": "ac"})
						st.Do(map[string]any{"out": "cc"}, nil, nil, nil)
					}
					t.Logf("%#v", trace)
				case tracing.ErrorTrace:
					t.Logf("%#v", trace)
				case flow.CeaseFlowTrace:
					return
				default:
					t.Logf("%#v", trace)
				}
			}
		}()

		select {
		case <-done:
		}

		t.Logf("%v\n", ins.Locator.CloneItems(data.LocatorObject))
		ins.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
