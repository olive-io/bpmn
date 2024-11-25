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

package process_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
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
