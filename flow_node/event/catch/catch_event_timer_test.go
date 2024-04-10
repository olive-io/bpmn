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

package catch_test

import (
	"context"
	"testing"
	"time"

	"github.com/olive-io/bpmn/schema"
	"github.com/stretchr/testify/require"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/event/catch"
	clock2 "github.com/olive-io/bpmn/pkg/clock"
	"github.com/olive-io/bpmn/pkg/timer"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/tracing"
)

var timerDoc schema.Definitions

func init() {
	LoadTestFile("testdata/intermediate_catch_event_timer.bpmn", &timerDoc)
}

func TestCatchEvent_Timer(t *testing.T) {
	processElement := (*timerDoc.Processes())[0]
	proc := process.New(&processElement, &timerDoc)
	fanOut := event.NewFanOut()
	c := clock2.NewMock()
	ctx := clock2.ToContext(context.Background(), c)
	tracer := tracing.NewTracer(ctx)
	eventInstanceBuilder := event.DefinitionInstanceBuildingChain(
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
