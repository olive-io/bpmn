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
	"fmt"
	"testing"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestBusinessRuleTask(t *testing.T) {
	var testTask schema.Definitions
	LoadTestFile("testdata/business_rule_task.bpmn", &testTask)

	engine := bpmn.NewEngine()
	ctx := context.TODO()
	if instance, err := engine.NewProcess(&testTask); err == nil {
		traces := instance.Tracer().Subscribe()
		err := instance.StartAll(ctx)
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
			case bpmn.FlowTrace:
			case bpmn.TaskTrace:
				var calledDecision *schema.ExtensionCalledDecision
				extension, _ := trace.GetActivity().Element().ExtensionElements()
				if extension != nil {
					calledDecision = extension.CalledDecision
				}
				if calledDecision != nil {
					//t.Logf("call decisionId [%s]", calledDecision.DecisionId)
					trace.Do(bpmn.DoWithResults(map[string]any{calledDecision.Result: "3"}))
				} else {
					trace.Do(bpmn.DoWithErr(fmt.Errorf("no result")))
				}
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
