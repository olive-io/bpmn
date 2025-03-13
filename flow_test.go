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

	"github.com/stretchr/testify/require"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/tracing"

	_ "github.com/olive-io/bpmn/v2/pkg/expression/expr"
)

func TestTrueFormalExpression(t *testing.T) {
	var testCondExpr schema.Definitions
	LoadTestFile("testdata/condexpr.bpmn", &testCondExpr)

	processElement := (*testCondExpr.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testCondExpr)
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
			case bpmn.CompletionTrace:
				if id, present := trace.Node.Id(); present {
					if *id == "end" {
						// success!
						break loop
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

func TestFalseFormalExpression(t *testing.T) {
	var testCondExprFalse schema.Definitions

	LoadTestFile("testdata/condexpr_false.bpmn", &testCondExprFalse)

	processElement := (*testCondExprFalse.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testCondExprFalse)
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
			case bpmn.CompletionTrace:
				if id, present := trace.Node.Id(); present {
					if *id == "end" {
						t.Fatalf("end should not have been reached")
					}
				}
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			case bpmn.CeaseFlowTrace:
				// success
				break loop
			default:
				//t.Logf("%#v", trace)
			}
		}
		instance.Tracer().Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}

func TestCondDataObject(t *testing.T) {
	var testCondDataObject schema.Definitions

	LoadTestFile("testdata/condexpr_dataobject.bpmn", &testCondDataObject)

	test := func(cond, expected string) func(t *testing.T) {
		return func(t *testing.T) {
			processElement := (*testCondDataObject.Processes())[0]
			proc := bpmn.NewProcess(&processElement, &testCondDataObject)
			if instance, err := proc.Instantiate(); err == nil {
				traces := instance.Tracer().Subscribe()
				// Set all data objects to false by default, except for `cond`
				for _, k := range []string{"cond1o", "cond2o"} {
					locator, _ := instance.Locator().FindIItemAwareLocator(data.LocatorObject)
					aware, found := locator.FindItemAwareByName(k)
					require.True(t, found)
					aware.Put(k == cond)
				}
				err := instance.StartAll()
				if err != nil {
					t.Fatalf("failed to run the instance: %s", err)
				}
			loop:
				for {
					trace := tracing.Unwrap(<-traces)
					switch trace := trace.(type) {
					case bpmn.VisitTrace:
						if id, present := trace.Node.Id(); present {
							if *id == expected {
								break loop
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
	}
	t.Run("cond1o", test("cond1o", "a1"))
	t.Run("cond2o", test("cond2o", "a2"))
}
