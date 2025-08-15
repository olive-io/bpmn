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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"

	_ "github.com/olive-io/bpmn/v2/pkg/expression/expr"
)

func TestParallelGateway(t *testing.T) {
	var testParallelGateway schema.Definitions
	LoadTestFile("testdata/parallel_gateway_fork_join.bpmn", &testParallelGateway)

	engine := bpmn.NewEngine()
	if instance, err := engine.NewProcess(&testParallelGateway); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		reached := make(map[string]int)
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.VisitTrace:
				//t.Logf("%#v", trace)
				if id, present := trace.Node.Id(); present {
					if counter, ok := reached[*id]; ok {
						reached[*id] = counter + 1
					} else {
						reached[*id] = 1
					}
				} else {
					t.Fatalf("can't find element with FlowNodeId %#v", id)
				}
			case bpmn.TaskTrace:
				trace.Do()
			case bpmn.CeaseFlowTrace:
				break loop
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				//t.Logf("%#v", trace)
			}
		}
		instance.Tracer().Unsubscribe(traces)

		assert.Equal(t, 1, reached["task1"])
		assert.Equal(t, 1, reached["task2"])
		assert.Equal(t, 2, reached["join"])
		assert.Equal(t, 1, reached["end"])
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestParallelGatewayMtoN(t *testing.T) {
	var testParallelGatewayMtoN schema.Definitions
	LoadTestFile("testdata/parallel_gateway_m_n.bpmn", &testParallelGatewayMtoN)

	engine := bpmn.NewEngine()
	if instance, err := engine.NewProcess(&testParallelGatewayMtoN); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		reached := make(map[string]int)
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.VisitTrace:
				//t.Logf("%#v", trace)
				if id, present := trace.Node.Id(); present {
					if counter, ok := reached[*id]; ok {
						reached[*id] = counter + 1
					} else {
						reached[*id] = 1
					}
				} else {
					t.Fatalf("can't find element with FlowNodeId %#v", id)
				}
			case bpmn.TaskTrace:
				trace.Do()
			case bpmn.CeaseFlowTrace:
				break loop
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				//t.Logf("%#v", trace)
			}
		}
		instance.Tracer().Unsubscribe(traces)

		assert.Equal(t, 3, reached["joinAndFork"])
		assert.Equal(t, 1, reached["task1"])
		assert.Equal(t, 1, reached["task2"])
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestParallelGatewayNtoM(t *testing.T) {
	var testParallelGatewayNtoM schema.Definitions
	LoadTestFile("testdata/parallel_gateway_n_m.bpmn", &testParallelGatewayNtoM)

	engine := bpmn.NewEngine()
	if instance, err := engine.NewProcess(&testParallelGatewayNtoM); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		reached := make(map[string]int)
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.VisitTrace:
				if id, present := trace.Node.Id(); present {
					if counter, ok := reached[*id]; ok {
						reached[*id] = counter + 1
					} else {
						reached[*id] = 1
					}
				} else {
					t.Fatalf("can't find element with FlowNodeId %#v", id)
				}
				//t.Logf("%#v", reached)
			case bpmn.TaskTrace:
				trace.Do()
			case bpmn.CeaseFlowTrace:
				break loop
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				//t.Logf("%#v", trace)
			}
		}
		instance.Tracer().Unsubscribe(traces)

		assert.Equal(t, 2, reached["joinAndFork"])
		assert.Equal(t, 1, reached["task1"])
		assert.Equal(t, 1, reached["task2"])
		assert.Equal(t, 1, reached["task3"])
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestParallelGatewayIncompleteJoin(t *testing.T) {
	var testParallelGatewayIncompleteJoin schema.Definitions
	LoadTestFile("testdata/parallel_gateway_fork_incomplete_join.bpmn", &testParallelGatewayIncompleteJoin)

	engine := bpmn.NewEngine()
	if instance, err := engine.NewProcess(&testParallelGatewayIncompleteJoin); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		reached := make(map[string]int)
	loop:
		for trace := range traces {
			trace = tracing.Unwrap(trace)
			switch trace := trace.(type) {
			case bpmn.IncomingFlowProcessedTrace:
				//t.Logf("%#v", trace)
				if nodeIdPtr, present := trace.Node.Id(); present {
					if *nodeIdPtr == "join" {
						source, err := trace.Flow.SequenceFlow().Source()
						assert.Nil(t, err)
						if idPtr, present := source.Id(); present {
							if *idPtr == "task1" {
								// task1 already came in and has been
								// processed
								break loop
							}
						}
					}
				}
			case bpmn.FlowTrace:
				if idPtr, present := trace.Source.Id(); present {
					if *idPtr == "join" {
						t.Fatalf("should not flow from join")
					}
				}
			case bpmn.VisitTrace:
				//t.Logf("%#v", trace)
				if id, present := trace.Node.Id(); present {
					if counter, ok := reached[*id]; ok {
						reached[*id] = counter + 1
					} else {
						reached[*id] = 1
					}
				} else {
					t.Fatalf("can't find element with FlowNodeId %#v", id)
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

		assert.Equal(t, 1, reached["task1"])
		assert.Equal(t, 0, reached["task2"])
		assert.Equal(t, 1, reached["join"])
		assert.Equal(t, 0, reached["end"])
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
