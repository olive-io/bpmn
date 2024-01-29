// Copyright 2023 The bpmn Authors
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

package catch_test

import (
	"context"
	"testing"
	"time"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/event/catch"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	clock2 "github.com/olive-io/bpmn/tools/clock"
	"github.com/olive-io/bpmn/tools/timer"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/require"
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
