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

package service_test

import (
	"context"
	"embed"
	"encoding/xml"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/olive-io/bpmn/schema"
	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
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
	LoadTestFile("testdata/service_task.bpmn", &testTask)
}

func TestServiceTask(t *testing.T) {
	processElement := (*testTask.Processes())[0]
	proc := process.New(&processElement, &testTask)
	options := []instance.Option{
		instance.WithVariables(map[string]any{
			"c": map[string]string{"name": "cc"},
			"d": []int32{1, 2, 3},
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
				case *activity.Trace:
					trace.Do()
					t.Logf("%#v", trace)
					properties := trace.GetProperties()
					t.Logf("properties: %#v", properties)
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
				case *activity.Trace:
					trace.Do(activity.WithErr(fmt.Errorf("text error")))
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
				case *activity.Trace:
					runnum += 1
					handler := make(chan flow_node.ErrHandler, 1)
					go func() {
						time.Sleep(time.Second * 1)
						retry := int32(1)
						handler <- flow_node.ErrHandler{Mode: flow_node.HandleRetry, Retries: retry}
					}()
					trace.Do(activity.WithErrHandle(fmt.Errorf("text error"), handler))
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
	locator := data.NewFlowDataLocator()
	options := []instance.Option{
		instance.WithLocator(locator),
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
				case *activity.Trace:
					assert.Equal(t, trace.GetDataObjects()["in"], map[string]any{"a": "ac"})
					trace.Do(activity.WithObjects(map[string]any{"out": map[string]any{"a": "cc"}}))
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

		t.Logf("%v\n", locator.CloneItems(data.LocatorObject))
		ins.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
