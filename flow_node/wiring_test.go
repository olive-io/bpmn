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

package flow_node_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"sync"
	"testing"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow_node"
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
