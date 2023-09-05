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

package inclusive_test

import (
	"context"
	"embed"
	"encoding/xml"
	"errors"
	"log"
	"testing"

	_ "github.com/olive-io/bpmn/expression/expr"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/gateway/inclusive"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/assert"
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

var testInclusiveGateway schema.Definitions

func init() {
	LoadTestFile("testdata/inclusive_gateway.bpmn", &testInclusiveGateway)
}

func TestInclusiveGateway(t *testing.T) {
	processElement := (*testInclusiveGateway.Processes())[0]
	proc := process.New(&processElement, &testInclusiveGateway)
	tracer := tracing.NewTracer(context.Background())
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 32))
	if inst, err := proc.Instantiate(instance.WithTracer(tracer)); err == nil {
		err := inst.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		endReached := 0
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case flow.Trace:
				for _, f := range trace.Flows {
					t.Logf("%#v", f.SequenceFlow())
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
			case flow.CeaseFlowTrace:
				// should only reach `end` once
				assert.Equal(t, 1, endReached)
				break loop
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				t.Logf("%#v", trace)
			}
		}
		inst.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

var testInclusiveGatewayDefault schema.Definitions

func init() {
	LoadTestFile("testdata/inclusive_gateway_default.bpmn", &testInclusiveGatewayDefault)
}

func TestInclusiveGatewayDefault(t *testing.T) {
	processElement := (*testInclusiveGatewayDefault.Processes())[0]
	proc := process.New(&processElement, &testInclusiveGatewayDefault)
	tracer := tracing.NewTracer(context.Background())
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 32))
	if inst, err := proc.Instantiate(instance.WithTracer(tracer)); err == nil {
		err := inst.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		endReached := 0
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case flow.Trace:
				for _, f := range trace.Flows {
					t.Logf("%#v", f.SequenceFlow())
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
			case flow.CeaseFlowTrace:
				// should only reach `end` once
				assert.Equal(t, 1, endReached)
				break loop
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				t.Logf("%#v", trace)
			}
		}
		inst.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

var testInclusiveGatewayNoDefault schema.Definitions

func init() {
	LoadTestFile("testdata/inclusive_gateway_no_default.bpmn", &testInclusiveGatewayNoDefault)
}

func TestInclusiveGatewayNoDefault(t *testing.T) {
	processElement := (*testInclusiveGatewayNoDefault.Processes())[0]
	proc := process.New(&processElement, &testInclusiveGatewayNoDefault)
	if inst, err := proc.Instantiate(); err == nil {
		traces := inst.Tracer.Subscribe()
		err := inst.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case tracing.ErrorTrace:
				var target inclusive.NoEffectiveSequenceFlows
				if errors.As(trace.Error, &target) {
					// success
					break loop
				} else {
					t.Fatalf("%#v", trace)
				}
			default:
				t.Logf("%#v", trace)
			}
		}
		inst.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
