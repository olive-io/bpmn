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

package flow_node_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow_node"
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

func TestNewWiring(t *testing.T) {
	var waitGroup sync.WaitGroup
	locator := data.NewFlowDataLocator()
	if proc, found := sampleDoc.FindBy(schema.ExactId("sample")); found {
		if flowNode, found := sampleDoc.FindBy(schema.ExactId("either")); found {
			node, err := flow_node.NewWiring(
				nil,
				proc.(*schema.Process),
				&defaultDefinitions,
				&flowNode.(*schema.ParallelGateway).FlowNode,
				event.VoidConsumer{},
				event.VoidSource{},
				tracing.NewTracer(context.Background()), flow_node.NewLockedFlowNodeMapping(),
				&waitGroup,
				event.WrappingDefinitionInstanceBuilder, locator,
			)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(node.Incoming))
			t.Logf("%+v", node.Incoming[0])
			if incomingSeqFlowId, present := node.Incoming[0].Id(); present {
				assert.Equal(t, *incomingSeqFlowId, "x1")
			} else {
				t.Fatalf("Sequence flow x1 has no matching ID")
			}
			assert.Equal(t, 2, len(node.Outgoing))
			if outgoingSeqFlowId, present := node.Outgoing[0].Id(); present {
				assert.Equal(t, *outgoingSeqFlowId, "x2")
			} else {
				t.Fatalf("Sequence flow x2 has no matching ID")
			}
			if outgoingSeqFlowId, present := node.Outgoing[1].Id(); present {
				assert.Equal(t, *outgoingSeqFlowId, "x3")
			} else {
				t.Fatalf("Sequence flow x3 has no matching ID")
			}
		} else {
			t.Fatalf("Can't find flow node `either`")
		}
	} else {
		t.Fatalf("Can't find process `sample`")
	}
}
