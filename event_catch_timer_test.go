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
	"time"

	"github.com/stretchr/testify/require"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/clock"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/timer"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestCatchEvent_Timer(t *testing.T) {
	var timerDoc schema.Definitions

	LoadTestFile("testdata/intermediate_catch_event_timer.bpmn", &timerDoc)

	engine := bpmn.NewEngine()
	fanOut := event.NewFanOut()
	c := clock.NewMock()
	ctx := clock.ToContext(context.Background(), c)
	tracer := tracing.NewTracer(ctx)
	eventInstanceBuilder := event.DefinitionInstanceBuildingChain(
		timer.EventDefinitionInstanceBuilder(ctx, fanOut, tracer),
	)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))

	proc, err := engine.NewProcess(&timerDoc, bpmn.WithTracer(tracer),
		bpmn.WithProcessEventDefinitionInstanceBuilder(eventInstanceBuilder),
		bpmn.WithEventEgress(fanOut),
		bpmn.WithEventIngress(fanOut))

	err = proc.StartAll()
	if err != nil {
		t.Fatalf("failed to run the instance: %s", err)
	}
	advancedTime := false
loop:
	for {
		trace := tracing.Unwrap(<-traces)
		switch trace := trace.(type) {
		case bpmn.ActiveListeningTrace:
			c.Add(1 * time.Minute)
			advancedTime = true
		case bpmn.CompletionTrace:
			if id, present := trace.Node.Id(); present {
				if *id == "end" {
					require.True(t, advancedTime)
					// success!
					break loop
				}

			}
		case bpmn.ErrorTrace:
			t.Fatalf("%#v", trace)
		default:
			//t.Logf("%#v", trace)
		}
	}
	proc.Tracer().Unsubscribe(traces)
}
