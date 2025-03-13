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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"

	_ "github.com/olive-io/bpmn/v2/pkg/expression/expr"
)

func TestExclusiveGateway(t *testing.T) {
	var testExclusiveGateway schema.Definitions
	LoadTestFile("testdata/exclusive_gateway.bpmn", &testExclusiveGateway)

	processElement := (*testExclusiveGateway.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testExclusiveGateway)
	if instance, err := proc.Instantiate(); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
				for _, f := range trace.Flows {
					//t.Logf("%#v", f.SequenceFlow())
					if target, err := f.SequenceFlow().Target(); err == nil {
						if id, present := target.Id(); present {
							assert.NotEqual(t, "task1", *id)
							if *id == "task2" {
								// reached task2 as expected
								break loop
							}
						} else {
							t.Fatalf("can't find target's FlowNodeId %#v", target)
						}

					} else {
						t.Fatalf("can't find sequence flow target: %#v", err)
					}
				}
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

func TestExclusiveGatewayWithDefault(t *testing.T) {
	var testExclusiveGatewayWithDefault schema.Definitions
	LoadTestFile("testdata/exclusive_gateway_default.bpmn", &testExclusiveGatewayWithDefault)

	processElement := (*testExclusiveGatewayWithDefault.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testExclusiveGatewayWithDefault)
	if instance, err := proc.Instantiate(); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
				for _, f := range trace.Flows {
					//t.Logf("%#v", f.SequenceFlow())
					if target, err := f.SequenceFlow().Target(); err == nil {
						if id, present := target.Id(); present {
							assert.NotEqual(t, "task1", *id)
							assert.NotEqual(t, "task2", *id)
							if *id == "default_task" {
								// reached default_task as expected
								break loop
							}
						} else {
							t.Fatalf("can't find target's FlowNodeId %#v", target)
						}

					} else {
						t.Fatalf("can't find sequence flow target: %#v", err)
					}
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
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestExclusiveGatewayWithNoDefault(t *testing.T) {
	var testExclusiveGatewayWithNoDefault schema.Definitions
	LoadTestFile("testdata/exclusive_gateway_no_default.bpmn", &testExclusiveGatewayWithNoDefault)

	processElement := (*testExclusiveGatewayWithNoDefault.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testExclusiveGatewayWithNoDefault)
	if instance, err := proc.Instantiate(); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
				for _, f := range trace.Flows {
					//t.Logf("%#v", f.SequenceFlow())
					if target, err := f.SequenceFlow().Target(); err == nil {
						if id, present := target.Id(); present {
							assert.NotEqual(t, "task1", *id)
							assert.NotEqual(t, "task2", *id)
						} else {
							t.Fatalf("can't find target's FlowNodeId %#v", target)
						}

					} else {
						t.Fatalf("can't find sequence flow target: %#v", err)
					}
				}
			case *bpmn.TaskTrace:
				trace.Do()
			case bpmn.ErrorTrace:
				var target bpmn.ExclusiveNoEffectiveSequenceFlows
				if errors.As(trace.Error, &target) {
					// success
					break loop
				} else {
					t.Fatalf("%#v", trace)
				}
			default:
				//t.Logf("%#v", trace)
			}
		}
		instance.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestExclusiveGatewayIncompleteJoin(t *testing.T) {
	var testExclusiveGatewayIncompleteJoin schema.Definitions

	LoadTestFile("testdata/exclusive_gateway_multiple_incoming.bpmn", &testExclusiveGatewayIncompleteJoin)

	processElement := (*testExclusiveGatewayIncompleteJoin.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testExclusiveGatewayIncompleteJoin)
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

		assert.Equal(t, 2, reached["exclusive"])
		assert.Equal(t, 2, reached["task2"])
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
