package catch_event

import (
	"context"
	"testing"
	"time"

	"github.com/olive-io/bpmn/clock"
	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/event/catch"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/test"
	"github.com/olive-io/bpmn/timer"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/require"
)

var timerDoc schema.Definitions

func init() {
	test.LoadTestFile("sample/catch_event/intermediate_catch_event_timer.bpmn", &timerDoc)
}

func TestCatchEvent_Timer(t *testing.T) {
	processElement := (*timerDoc.Processes())[0]
	proc := process.New(&processElement, &timerDoc)
	fanOut := event.NewFanOut()
	c := clock.NewMock()
	ctx := clock.ToContext(context.Background(), c)
	tracer := tracing.NewTracer(ctx)
	eventInstanceBuilder := event.IDefinitionInstanceBuildingChain(
		timer.EventDefinitionInstanceBuilder(ctx, fanOut, tracer),
	)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
	if i, err := proc.Instantiate(
		instance.WithTracer(tracer),
		instance.WithEventDefinitionInstanceBuilder(eventInstanceBuilder),
		instance.WithEventEgress(fanOut),
		instance.WithEventIngress(fanOut),
	); err == nil {
		err := i.StartAll(ctx)
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
		advancedTime := false
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case catch.ActiveListeningTrace:
				c.Add(1 * time.Minute)
				advancedTime = true
			case flow.CompletionTrace:
				if id, present := trace.Node.Id(); present {
					if *id == "end" {
						require.True(t, advancedTime)
						// success!
						break loop
					}

				}
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				t.Logf("%#v", trace)
			}
		}
		i.Tracer.Unsubscribe(traces)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
