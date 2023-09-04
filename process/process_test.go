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

package process_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/olive-io/bpmn/flow_node"
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

var defaultDefinitions = schema.DefaultDefinitions()

var sampleDoc schema.Definitions

func init() {
	LoadTestFile("testdata/sample.bpmn", &sampleDoc)
}

func TestExplicitInstantiation(t *testing.T) {
	if proc, found := sampleDoc.FindBy(schema.ExactId("sample")); found {
		process := process.New(proc.(*schema.Process), &defaultDefinitions)
		inst, err := process.Instantiate()
		assert.Nil(t, err)
		assert.NotNil(t, inst)
	} else {
		t.Fatalf("Can't find process `sample`")
	}
}

func TestCancellation(t *testing.T) {
	if proc, found := sampleDoc.FindBy(schema.ExactId("sample")); found {
		ctx, cancel := context.WithCancel(context.Background())

		process := process.New(proc.(*schema.Process), &defaultDefinitions, process.WithContext(ctx))

		tracer := tracing.NewTracer(ctx)
		traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))

		inst, err := process.Instantiate(instance.WithContext(ctx), instance.WithTracer(tracer))
		assert.Nil(t, err)
		assert.NotNil(t, inst)

		cancel()

		cancelledFlowNodes := make([]schema.FlowNodeInterface, 0)

		for trace := range traces {
			trace = tracing.Unwrap(trace)
			switch trace := trace.(type) {
			case flow_node.CancellationTrace:
				cancelledFlowNodes = append(cancelledFlowNodes, trace.Node)
			default:
			}
		}

		assert.NotEmpty(t, cancelledFlowNodes)
	} else {
		t.Fatalf("Can't find process `sample`")
	}
}
