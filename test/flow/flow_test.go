package flow

import (
	"context"
	"testing"

	_ "github.com/olive-io/bpmn/expression/expr"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/test"
	"github.com/olive-io/bpmn/tracing"
)

var testCondExpr schema.Definitions

func init() {
	test.LoadTestFile("sample/flow/condexpr.bpmn", &testCondExpr)
}

func TestTrueFormalExpression(t *testing.T) {
	processElement := (*testCondExpr.Processes())[0]
	proc := process.New(&processElement, &testCondExpr)
	if instance, err := proc.Instantiate(); err == nil {
		traces := instance.Tracer.Subscribe()
		err := instance.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case flow.CompletionTrace:
				if id, present := trace.Node.Id(); present {
					if *id == "end" {
						// success!
						break loop
					}

				}
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				t.Logf("%#v", trace)
			}
		}
		instance.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

var testCondExprFalse schema.Definitions

func init() {
	test.LoadTestFile("sample/flow/condexpr_false.bpmn", &testCondExprFalse)
}

func TestFalseFormalExpression(t *testing.T) {
	processElement := (*testCondExprFalse.Processes())[0]
	proc := process.New(&processElement, &testCondExprFalse)
	if instance, err := proc.Instantiate(); err == nil {
		traces := instance.Tracer.Subscribe()
		err := instance.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case flow.CompletionTrace:
				if id, present := trace.Node.Id(); present {
					if *id == "end" {
						t.Fatalf("end should not have been reached")
					}
				}
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			case flow.CeaseFlowTrace:
				// success
				break loop
			default:
				t.Logf("%#v", trace)
			}
		}
		instance.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

var testCondDataObject schema.Definitions

func init() {
	test.LoadTestFile("sample/flow/condexpr_dataobject.bpmn", &testCondDataObject)
}

func TestCondDataObject(t *testing.T) {
	test := func(cond, expected string) func(t *testing.T) {
		return func(t *testing.T) {
			processElement := (*testCondDataObject.Processes())[0]
			proc := process.New(&processElement, &testCondDataObject)
			if instance, err := proc.Instantiate(); err == nil {
				traces := instance.Tracer.Subscribe()
				// Set all data objects to false by default, except for `cond`
				for _, k := range []string{"cond1o", "cond2o"} {
					_ = k
					//itemAware, found := instance.FindItemAwareByName(k)
					//require.True(t, found)
					//itemAware.Put(context.Background(), k == cond)
				}
				err := instance.StartAll(context.Background())
				if err != nil {
					t.Fatalf("failed to run the instance: %s", err)
				}
			loop:
				for {
					trace := tracing.Unwrap(<-traces)
					switch trace := trace.(type) {
					case flow.VisitTrace:
						if id, present := trace.Node.Id(); present {
							if *id == expected {
								break loop
							}
						}
					case tracing.ErrorTrace:
						t.Fatalf("%#v", trace)
					default:
						t.Logf("%#v", trace)
					}
				}
				instance.Tracer.Unsubscribe(traces)
			} else {
				t.Fatalf("failed to instantiate the process: %s", err)
			}
		}
	}
	t.Run("cond1o", test("cond1o", "a1"))
	t.Run("cond2o", test("cond2o", "a2"))
}
