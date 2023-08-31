package activity

import (
	"context"
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
	"github.com/olive-io/bpmn/test"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/assert"

	_ "github.com/stretchr/testify/assert"
)

var testDoc schema.Definitions

func init() {
	test.LoadTestFile("sample/activity/boundary_event.bpmn", &testDoc)
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

func testBoundaryEvent(t *testing.T, boundary string, test func(visited map[string]bool), events ...event.Event) {
	processElement := (*testDoc.Processes())[0]
	proc := process.New(&processElement, &testDoc)
	ready := make(chan bool)

	// explicit tracer
	tracer := tracing.NewTracer(context.Background())
	// this gives us some room when instance starts up
	traces := tracer.SubscribeChannel(make(chan tracing.Trace, 32))

	if inst, err := proc.Instantiate(instance.WithTracer(tracer)); err == nil {
		if node, found := testDoc.FindBy(schema.ExactId("task")); found {
			if taskNode, found := inst.FlowNodeMapping().
				ResolveElementToFlowNode(node.(schema.FlowNodeInterface)); found {
				harness := taskNode.(*activity.Harness)
				aTask := harness.Activity().(*task.Task)
				aTask.SetBody(func(task *task.Task, ctx context.Context) flow_node.Action {
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

		err := inst.StartAll(context.Background())
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
