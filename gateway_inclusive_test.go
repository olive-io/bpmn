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
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"

	_ "github.com/olive-io/bpmn/v2/pkg/expression/expr"
)

func TestInclusiveGateway(t *testing.T) {
	var testInclusiveGateway schema.Definitions
	LoadTestFile("testdata/inclusive_gateway.bpmn", &testInclusiveGateway)

	processElement := (*testInclusiveGateway.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testInclusiveGateway)
	tracer := tracing.NewTracer(context.Background())
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 32))
	if inst, err := proc.Instantiate(bpmn.WithTracer(tracer)); err == nil {
		err := inst.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		endReached := 0
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
				for _, f := range trace.Flows {
					if target, err := f.SequenceFlow().Target(); err == nil {
						if id, present := target.Id(); present {
							assert.NotEqual(t, "a3", *id)
							if *id == "end" {
								// reached end
								endReached++
								continue
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
			case bpmn.CeaseFlowTrace:
				// should only reach `end` once
				assert.Equal(t, 1, endReached)
				break loop
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				//t.Logf("%#v", trace)
			}
		}
		inst.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestInclusiveGatewayDefault(t *testing.T) {
	var testInclusiveGatewayDefault schema.Definitions
	LoadTestFile("testdata/inclusive_gateway_default.bpmn", &testInclusiveGatewayDefault)

	processElement := (*testInclusiveGatewayDefault.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testInclusiveGatewayDefault)
	tracer := tracing.NewTracer(context.Background())
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 32))
	if inst, err := proc.Instantiate(bpmn.WithTracer(tracer)); err == nil {
		err := inst.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		endReached := 0
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
				for _, f := range trace.Flows {
					//t.Logf("%#v", f.SequenceFlow())
					if target, err := f.SequenceFlow().Target(); err == nil {
						if id, present := target.Id(); present {
							assert.NotEqual(t, "a1", *id)
							assert.NotEqual(t, "a2", *id)
							if *id == "end" {
								// reached end
								endReached++
								continue
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
			case bpmn.CeaseFlowTrace:
				// should only reach `end` once
				assert.Equal(t, 1, endReached)
				break loop
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				//t.Logf("%#v", trace)
			}
		}
		inst.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestInclusiveGatewayNoDefault(t *testing.T) {
	var testInclusiveGatewayNoDefault schema.Definitions
	LoadTestFile("testdata/inclusive_gateway_no_default.bpmn", &testInclusiveGatewayNoDefault)

	processElement := (*testInclusiveGatewayNoDefault.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testInclusiveGatewayNoDefault)
	if inst, err := proc.Instantiate(); err == nil {
		traces := inst.Tracer().Subscribe()
		err := inst.StartAll()
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.ErrorTrace:
				var target bpmn.InclusiveNoEffectiveSequenceFlows
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
		inst.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
