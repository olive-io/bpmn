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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestEventThrow(t *testing.T) {
	var testTask schema.Definitions

	LoadTestFile("testdata/intermediate_event_throw.bpmn", &testTask)

	engine := bpmn.NewEngine()
	ctx := context.Background()
	visited := make([]string, 0)
	instance, err := engine.NewProcess(&testTask)
	if err != nil {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
	traces := instance.Tracer().Subscribe()
	defer instance.Tracer().Unsubscribe(traces)
	err = instance.StartAll(ctx)
	if err != nil {
		t.Fatalf("failed to run the instance: %s", err)
	}

loop:
	for {
		trace := tracing.Unwrap(<-traces)
		switch trace := trace.(type) {
		case bpmn.FlowTrace:

		case bpmn.TaskTrace:
			name, ok := trace.GetActivity().Element().Name()
			if ok {
				visited = append(visited, *name)
			}
			trace.Do()
		case bpmn.ActiveListeningTrace:
			for _, eventDefinition := range trace.Node.SignalEventDefinitionField {
				ref, ok := eventDefinition.SignalRef()
				if ok {
					instance.ConsumeEvent(event.NewSignalEvent(string(*ref)))
				}
			}
			for _, eventDefinition := range trace.Node.MessageEventDefinitionField {
				ref, ok := eventDefinition.MessageRef()
				if ok {
					instance.ConsumeEvent(event.NewMessageEvent(string(*ref), (*string)(eventDefinition.OperationRefField)))
				}
			}
		case bpmn.CeaseFlowTrace:
			break loop
		case bpmn.TerminationTrace:
		//	t.Logf("%#v", trace)
		case bpmn.ErrorTrace:
			t.Fatalf("%#v", trace)
		default:
			//t.Logf("%#v", trace)
		}
	}

	sort.Strings(visited)
	assert.Equal(t, visited, []string{"s1", "s2", "t1", "t2"})
}
