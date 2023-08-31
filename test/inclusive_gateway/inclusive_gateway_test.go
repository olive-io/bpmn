package inclusive_gateway

import (
	"context"
	"errors"
	"testing"

	_ "github.com/olive-io/bpmn/expression/expr"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/gateway/inclusive"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/test"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/assert"
)

var testInclusiveGateway schema.Definitions

func init() {
	test.LoadTestFile("sample/inclusive_gateway/inclusive_gateway.bpmn", &testInclusiveGateway)
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
	test.LoadTestFile("sample/inclusive_gateway/inclusive_gateway_default.bpmn", &testInclusiveGatewayDefault)
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
	test.LoadTestFile("sample/inclusive_gateway/inclusive_gateway_no_default.bpmn", &testInclusiveGatewayNoDefault)
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
