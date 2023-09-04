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

package activity_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/flow_node/activity/task"
	"github.com/olive-io/bpmn/flow_node/event/catch"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/assert"

	_ "github.com/stretchr/testify/assert"
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

var testDoc schema.Definitions

func init() {
	LoadTestFile("testdata/boundary_event.bpmn", &testDoc)
}

func TestInterruptingEvent(t *testing.T) {
	testBoundaryEvent(t, "sig1listener", func(visited map[string]bool) {
		assert.False(t, visited["uninterrupted"])
		assert.True(t, visited["interrupted"])
		assert.True(t, visited["end"])
	}, event.NewSignalEvent("sig1"))
}

func TestNonInterruptingEvent(t *testing.T) {
	testBoundaryEvent(t, "sig2listener", func(visited map[string]bool) {
		assert.False(t, visited["interrupted"])
		assert.True(t, visited["uninterrupted"])
		assert.True(t, visited["end"])
	}, event.NewSignalEvent("sig2"))
}

func testBoundaryEvent(t *testing.T, boundary string, test func(visited map[string]bool), events ...event.IEvent) {
	processElement := (*testDoc.Processes())[0]
	proc := process.New(&processElement, &testDoc)
	ready := make(chan bool)

	// explicit tracer
	tracer := tracing.NewTracer(context.Background())
	// this gives us some room when instance starts up
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 32))

	if inst, err := proc.Instantiate(instance.WithTracer(tracer)); err == nil {
		if node, found := testDoc.FindBy(schema.ExactId("task")); found {
			if taskNode, found := inst.FlowNodeMapping().
				ResolveElementToFlowNode(node.(schema.FlowNodeInterface)); found {
				harness := taskNode.(*activity.Harness)
				aTask := harness.Activity().(*task.Task)
				aTask.SetBody(func(task *task.Task, ctx context.Context) flow_node.IAction {
					select {
					case <-ready:
						return flow_node.FlowAction{SequenceFlows: flow_node.AllSequenceFlows(&task.Wiring.Outgoing)}
					case <-ctx.Done():
						return flow_node.CompleteAction{}
					}
				})
			} else {
				t.Fatalf("failed to get the flow node `task`")
			}
		} else {
			t.Fatalf("failed to get the flow node element for `task`")
		}

		err = inst.StartAll(context.Background(), nil)
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}

		listening := false
		activeBoundary := false

		for {
			if listening && activeBoundary {
				break
			}
			trace := tracing.Unwrap(<-traces)
			t.Logf("%#v", trace)
			switch trace := trace.(type) {
			case catch.ActiveListeningTrace:
				// Ensure the boundary event listener is actually listening
				if id, present := trace.Node.Id(); present {
					if *id == boundary {
						// it is indeed listening
						listening = true
					}

				}
			case activity.ActiveBoundaryTrace:
				if id, present := trace.Node.Id(); present && trace.Start {
					if *id == "task" {
						// task has reached its active boundary
						activeBoundary = true
					}
				}
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				t.Logf("%#v", trace)
			}
		}

		for _, evt := range events {
			_, err = inst.ConsumeEvent(evt)
			assert.Nil(t, err)
		}
		visited := make(map[string]bool)
	loop1:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case flow.VisitTrace:
				if id, present := trace.Node.Id(); present {
					if *id == "uninterrupted" {
						// we're here to we can release the task
						ready <- true
					}
					visited[*id] = true
					if *id == "end" {
						break loop1
					}
				}
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				t.Logf("%#v", trace)
			}
		}
		inst.Tracer.Unsubscribe(traces)

		test(visited)

	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}