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

	engine := bpmn.NewEngine()
	ctx := context.Background()
	tracer := tracing.NewTracer(ctx)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 32))
	if inst, err := engine.NewProcess(&testInclusiveGateway, bpmn.WithTracer(tracer)); err == nil {
		err := inst.StartAll(ctx)
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		endReached := 0
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
			case bpmn.TaskTrace:
				trace.Do()
			case bpmn.CeaseFlowTrace:
				// should only reach `end` once
				assert.Equal(t, 1, endReached)
				break loop
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
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

	engine := bpmn.NewEngine()
	ctx := context.Background()
	tracer := tracing.NewTracer(ctx)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 32))
	if inst, err := engine.NewProcess(&testInclusiveGatewayDefault, bpmn.WithTracer(tracer)); err == nil {
		err := inst.StartAll(ctx)
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		endReached := 0
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
			case bpmn.TaskTrace:
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

	engine := bpmn.NewEngine()
	ctx := context.Background()
	if inst, err := engine.NewProcess(&testInclusiveGatewayNoDefault); err == nil {
		traces := inst.Tracer().Subscribe()
		err := inst.StartAll(ctx)
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
