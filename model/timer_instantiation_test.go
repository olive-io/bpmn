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

package model_test

import (
	"context"
	"testing"
	"time"

	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/model"
	"github.com/olive-io/bpmn/schema"
	clock2 "github.com/olive-io/bpmn/tools/clock"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/require"
)

var testTimerStartEventInstantiation schema.Definitions

func init() {
	LoadTestFile("testdata/instantiate_timer_start_event.bpmn", &testTimerStartEventInstantiation)
}

func TestTimerStartEventInstantiation(t *testing.T) {
	c := clock2.NewMock()
	ctx := clock2.ToContext(context.Background(), c)
	tracer := tracing.NewTracer(ctx)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
	m := model.New(&testTimerStartEventInstantiation, model.WithContext(ctx), model.WithTracer(tracer))
	err := m.Run(ctx)
	require.Nil(t, err)
loop:
	for {
		select {
		case trace := <-traces:
			_, ok := trace.(flow.Trace)
			// Should not flow
			require.False(t, ok)
		default:
			break loop
		}
	}
	// Advance clock by 1M
	c.Add(1 * time.Minute)
loop1:
	for {
		trace := tracing.Unwrap(<-traces)
		switch trace := trace.(type) {
		case flow.VisitTrace:
			if idPtr, present := trace.Node.Id(); present {
				if *idPtr == "end" {
					// we've reached the desired outcome
					break loop1
				}
			}
		default:
			t.Logf("%#v", trace)
		}
	}
}

var testRecurringTimerStartEventInstantiation schema.Definitions

func init() {
	LoadTestFile("testdata/instantiate_recurring_timer_start_event.bpmn", &testRecurringTimerStartEventInstantiation)
}

func TestRecurringTimerStartEventInstantiation(t *testing.T) {
	c := clock2.NewMock()
	ctx := clock2.ToContext(context.Background(), c)
	tracer := tracing.NewTracer(ctx)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
	m := model.New(&testRecurringTimerStartEventInstantiation, model.WithContext(ctx), model.WithTracer(tracer))
	err := m.Run(ctx)
	require.Nil(t, err)
loop:
	for {
		select {
		case trace := <-traces:
			trace = tracing.Unwrap(trace)
			_, ok := trace.(flow.Trace)
			// Should not flow
			require.False(t, ok)
		default:
			break loop
		}
	}
	// Test for some arbitrary number of recurrences (say, 10?)
	for i := 0; i < 10; i++ {
		// Advance clock by 1M
		c.Add(1 * time.Minute)
	loop1:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case flow.VisitTrace:
				if idPtr, present := trace.Node.Id(); present {
					if *idPtr == "end" {
						// we've reached the desired outcome
						break loop1
					}
				}
			default:
				t.Logf("%#v", trace)
			}
		}
	}
}
