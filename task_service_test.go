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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestServiceTask(t *testing.T) {
	var testTask schema.Definitions
	LoadTestFile("testdata/service_task.bpmn", &testTask)

	processElement := (*testTask.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testTask)
	options := []bpmn.Option{
		bpmn.WithVariables(map[string]any{
			"c": map[string]string{"name": "cc"},
			"d": []int32{1, 2, 3},
		}),
		bpmn.WithDataObjects(map[string]any{
			"a": struct{}{},
		}),
	}
	if ins, err := proc.Instantiate(options...); err == nil {
		traces := ins.Tracer().Subscribe()
		err = ins.StartAll()
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
				case bpmn.FlowTrace:
				case bpmn.TaskTrace:
					trace.Do(bpmn.DoWithResults(map[string]interface{}{"foo": 1, "bar": "2"}))
					//t.Logf("%#v", trace)
					//properties := trace.GetProperties()
					//t.Logf("properties: %#v", properties)
				case bpmn.ErrorTrace:
					t.Errorf("%#v", trace)
					return
				case bpmn.CeaseFlowTrace:
					return
				default:
					//t.Logf("%#v", trace)
				}
			}
		}()

		select {
		case <-done:
		}

		foo, _ := ins.Locator().GetVariable("foo")
		assert.Equal(t, foo, 1)
		_, ok := ins.Locator().GetVariable("bar")
		assert.True(t, !ok)
		ins.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestServiceTaskWithError(t *testing.T) {
	var testTask schema.Definitions
	LoadTestFile("testdata/service_task.bpmn", &testTask)

	processElement := (*testTask.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testTask)
	var te error
	if ins, err := proc.Instantiate(); err == nil {
		traces := ins.Tracer().Subscribe()
		err := ins.StartAll()
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
				case bpmn.FlowTrace:
				case bpmn.TaskTrace:
					trace.Do(bpmn.DoWithErr(fmt.Errorf("text error")))
					//t.Logf("%#v", trace)
				case bpmn.ErrorTrace:
					te = trace.Error
					//t.Logf("%#v", trace)
				case bpmn.CeaseFlowTrace:
					return
				default:
					//t.Logf("%#v", trace)
				}
			}
		}()

		select {
		case <-done:
		}

		ins.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}

	assert.Equal(t, te.(errors.TaskExecError).Reason, "text error")
}

func TestServiceTaskWithRetry(t *testing.T) {
	var testTask schema.Definitions
	LoadTestFile("testdata/service_task.bpmn", &testTask)

	processElement := (*testTask.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testTask)
	var te error
	runnum := 0
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	if ins, err := proc.Instantiate(); err == nil {
		traces := ins.Tracer().Subscribe()
		err = ins.StartAll()
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
				case bpmn.FlowTrace:
				case bpmn.TaskTrace:
					runnum += 1
					handler := make(chan bpmn.ErrHandler, 1)
					go func() {
						time.Sleep(time.Second * 1)
						retry := int32(1)
						handler <- bpmn.ErrHandler{Mode: bpmn.RetryMode, Retries: retry}
					}()
					trace.Do(bpmn.DoWithErrHandle(fmt.Errorf("text error"), handler))
					//t.Logf("%#v", trace)
				case bpmn.ErrorTrace:
					te = trace.Error
					//t.Logf("%#v", trace)
				case bpmn.CeaseFlowTrace:
					return
				default:
					// t.Logf("%#v", trace)
				}
			}
		}()

		select {
		case <-done:
		}

		ins.Tracer().Unsubscribe(traces)
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
	locator := data.NewFlowDataLocator()
	options := []bpmn.Option{
		bpmn.WithLocator(locator),
		bpmn.WithDataObjects(map[string]any{"DataObject_0yhrl3s": map[string]any{"a": "ac"}}),
	}
	proc := bpmn.NewProcess(&processElement, task)
	if ins, err := proc.Instantiate(options...); err == nil {
		traces := ins.Tracer().Subscribe()
		err := ins.StartAll()
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
				case bpmn.FlowTrace:
				case bpmn.TaskTrace:
					assert.Equal(t, trace.GetDataObjects()["in"], map[string]any{"a": "ac"})
					trace.Do(bpmn.DoWithResults(map[string]any{"out": map[string]any{"a": "cc"}}))
					//t.Logf("%#v", trace)
				case bpmn.ErrorTrace:
					t.Logf("%#v", trace)
				case bpmn.CeaseFlowTrace:
					return
				default:
					//t.Logf("%#v", trace)
				}
			}
		}()

		select {
		case <-done:
		}

		t.Logf("%v\n", locator.CloneItems(data.LocatorObject))
		ins.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
