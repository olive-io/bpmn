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

package bpmn

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

//go:embed testdata/sample.bpmn
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

func TestNewWiring(t *testing.T) {
	var defaultDefinitions = schema.DefaultDefinitions()

	var sampleDoc schema.Definitions
	LoadTestFile("testdata/sample.bpmn", &sampleDoc)

	var waitGroup sync.WaitGroup
	locator := data.NewFlowDataLocator()
	if proc, found := sampleDoc.FindBy(schema.ExactId("sample")); found {
		if flowNode, found := sampleDoc.FindBy(schema.ExactId("either")); found {
			node, err := newWiring(
				nil,
				proc.(*schema.Process),
				&defaultDefinitions,
				&flowNode.(*schema.ParallelGateway).FlowNode,
				event.VoidConsumer{},
				event.VoidSource{},
				tracing.NewTracer(context.Background()), NewLockedFlowNodeMapping(),
				&waitGroup,
				event.WrappingDefinitionInstanceBuilder,
				locator,
			)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(node.incoming))
			//t.Logf("%+v", node.Incoming[0])
			if incomingSeqFlowId, present := node.incoming[0].Id(); present {
				assert.Equal(t, *incomingSeqFlowId, "x1")
			} else {
				t.Fatalf("Sequence flow x1 has no matching ID")
			}
			assert.Equal(t, 2, len(node.outgoing))
			if outgoingSeqFlowId, present := node.outgoing[0].Id(); present {
				assert.Equal(t, *outgoingSeqFlowId, "x2")
			} else {
				t.Fatalf("Sequence flow x2 has no matching ID")
			}
			if outgoingSeqFlowId, present := node.outgoing[1].Id(); present {
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
