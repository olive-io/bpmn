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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestExplicitInstantiation(t *testing.T) {
	var defaultDefinitions = schema.DefaultDefinitions()
	var sampleDoc schema.Definitions
	LoadTestFile("testdata/sample.bpmn", &sampleDoc)

	if proc, found := sampleDoc.FindBy(schema.ExactId("sample")); found {
		process := bpmn.NewProcess(proc.(*schema.Process), &defaultDefinitions)
		inst, err := process.Instantiate()
		assert.Nil(t, err)
		assert.NotNil(t, inst)
	} else {
		t.Fatalf("Can't find process `sample`")
	}
}

func TestCancellation(t *testing.T) {
	var defaultDefinitions = schema.DefaultDefinitions()
	var sampleDoc schema.Definitions
	LoadTestFile("testdata/sample.bpmn", &sampleDoc)

	if proc, found := sampleDoc.FindBy(schema.ExactId("sample")); found {
		ctx, cancel := context.WithCancel(context.Background())

		proc := bpmn.NewProcess(proc.(*schema.Process), &defaultDefinitions, bpmn.WithContext(ctx))

		tracer := tracing.NewTracer(ctx)
		traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))

		inst, err := proc.Instantiate(bpmn.WithContext(ctx), bpmn.WithTracer(tracer))
		assert.Nil(t, err)
		assert.NotNil(t, inst)

		cancel()

		cancelledFlowNodes := make([]schema.FlowNodeInterface, 0)

		for trace := range traces {
			trace = tracing.Unwrap(trace)
			switch trace := trace.(type) {
			case bpmn.CancellationFlowNodeTrace:
				cancelledFlowNodes = append(cancelledFlowNodes, trace.Node)
			default:
			}
		}

		assert.NotEmpty(t, cancelledFlowNodes)
	} else {
		t.Fatalf("Can't find process `sample`")
	}
}
