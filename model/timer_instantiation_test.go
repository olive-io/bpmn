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

package model_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/model"
	"github.com/olive-io/bpmn/v2/pkg/clock"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

var testTimerStartEventInstantiation schema.Definitions

func init() {
	LoadTestFile("testdata/instantiate_timer_start_event.bpmn", &testTimerStartEventInstantiation)
}

func TestTimerStartEventInstantiation(t *testing.T) {
	c := clock.NewMock()
	ctx := clock.ToContext(context.Background(), c)
	tracer := tracing.NewTracer(ctx)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
	m := model.New(&testTimerStartEventInstantiation, model.WithContext(ctx), model.WithTracer(tracer))
	err := m.Run(ctx)
	require.Nil(t, err)
loop:
	for {
		select {
		case trace := <-traces:
			_, ok := trace.(bpmn.FlowTrace)
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
		case bpmn.VisitTrace:
			if idPtr, present := trace.Node.Id(); present {
				if *idPtr == "end" {
					// we've reached the desired outcome
					break loop1
				}
			}
		default:
			//t.Logf("%#v", trace)
		}
	}
}

var testRecurringTimerStartEventInstantiation schema.Definitions

func init() {
	LoadTestFile("testdata/instantiate_recurring_timer_start_event.bpmn", &testRecurringTimerStartEventInstantiation)
}

func TestRecurringTimerStartEventInstantiation(t *testing.T) {
	c := clock.NewMock()
	ctx := clock.ToContext(context.Background(), c)
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
			_, ok := trace.(bpmn.FlowTrace)
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
			case bpmn.VisitTrace:
				if idPtr, present := trace.Node.Id(); present {
					if *idPtr == "end" {
						// we've reached the desired outcome
						break loop1
					}
				}
			default:
				//t.Logf("%#v", trace)
			}
		}
	}
}
