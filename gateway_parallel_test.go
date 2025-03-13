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

	_ "github.com/olive-io/bpmn/v2/pkg/expression/expr"
)

func TestParallelGateway(t *testing.T) {
	var testParallelGateway schema.Definitions
	LoadTestFile("testdata/parallel_gateway_fork_join.bpmn", &testParallelGateway)

	processElement := (*testParallelGateway.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testParallelGateway)
	if instance, err := proc.Instantiate(); err == nil {
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
			case *bpmn.TaskTrace:
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

	processElement := (*testParallelGatewayMtoN.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testParallelGatewayMtoN)
	if instance, err := proc.Instantiate(); err == nil {
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
			case *bpmn.TaskTrace:
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

	processElement := (*testParallelGatewayNtoM.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testParallelGatewayNtoM)
	if instance, err := proc.Instantiate(); err == nil {
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
			case *bpmn.TaskTrace:
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

	processElement := (*testParallelGatewayIncompleteJoin.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testParallelGatewayIncompleteJoin)
	if instance, err := proc.Instantiate(); err == nil {
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
			case *bpmn.TaskTrace:
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
