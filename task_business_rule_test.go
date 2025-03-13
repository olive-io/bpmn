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

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestBusinessRuleTask(t *testing.T) {
	var testTask schema.Definitions
	LoadTestFile("testdata/business_rule_task.bpmn", &testTask)

	processElement := (*testTask.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &testTask)
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
			case *bpmn.TaskTrace:
				calledDecision := trace.Context().Value(bpmn.BusinessRuleTaskKey{}).(*schema.ExtensionCalledDecision)
				//t.Logf("call decisionId [%s]", calledDecision.DecisionId)
				trace.Do(bpmn.DoWithResults(map[string]any{calledDecision.Result: "3"}))
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			case bpmn.CeaseFlowTrace:
				break loop
			default:
				//t.Logf("%#v", trace)
			}
		}
		instance.Tracer().Unsubscribe(traces)
		//t.Logf("%#v", instance.Locator().CloneVariables())
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
