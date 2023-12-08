// Copyright 2023 The olive Authors
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

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/model"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/require"
)

var testStartEventInstantiation schema.Definitions

func init() {
	LoadTestFile("testdata/instantiate_start_event.bpmn", &testStartEventInstantiation)
}

func TestStartEventInstantiation(t *testing.T) {
	ctx := context.Background()
	tracer := tracing.NewTracer(ctx)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
	m := model.New(&testStartEventInstantiation, model.WithContext(ctx), model.WithTracer(tracer))
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
	// Test simple event instantiation
	_, err = m.ConsumeEvent(event.NewSignalEvent("sig1"))
	require.Nil(t, err)
loop1:
	for {
		trace := tracing.Unwrap(<-traces)
		switch trace := trace.(type) {
		case flow.VisitTrace:
			if idPtr, present := trace.Node.Id(); present {
				if *idPtr == "sig1a" {
					// we've reached the desired outcome
					break loop1
				}
			}
		default:
			t.Logf("%#v", trace)
		}
	}
}

func TestMultipleStartEventInstantiation(t *testing.T) {
	ctx := context.Background()
	tracer := tracing.NewTracer(ctx)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
	m := model.New(&testStartEventInstantiation, model.WithContext(ctx), model.WithTracer(tracer))
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
	// Test multiple event instantiation
	_, err = m.ConsumeEvent(event.NewSignalEvent("sig2"))
	require.Nil(t, err)
loop1:
	for {
		trace := tracing.Unwrap(<-traces)
		switch trace := trace.(type) {
		case flow.VisitTrace:
			if idPtr, present := trace.Node.Id(); present {
				if *idPtr == "sig2_sig3a" {
					// we've reached the desired outcome
					break loop1
				}
			}
		default:
			t.Logf("%#v", trace)
		}
	}
}

func TestParallelMultipleStartEventInstantiation(t *testing.T) {
	ctx := context.Background()
	tracer := tracing.NewTracer(ctx)
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
	m := model.New(&testStartEventInstantiation, model.WithContext(ctx), model.WithTracer(tracer))
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
	// Test parallel multiple event instantiation
	signalEvent := event.NewSignalEvent("sig2")
	_, err = m.ConsumeEvent(signalEvent)
	require.Nil(t, err)
	sig3sent := false
loop1:
	for {
		trace := tracing.Unwrap(<-traces)
		switch trace := trace.(type) {
		case model.EventInstantiationAttemptedTrace:
			if !sig3sent && signalEvent == trace.Event {
				if idPtr, present := trace.Element.Id(); present && *idPtr == "ParallelMultipleStartEvent" {
					_, err = m.ConsumeEvent(event.NewSignalEvent("sig3"))
					require.Nil(t, err)
					sig3sent = true
				}
			}
		case flow.VisitTrace:
			if idPtr, present := trace.Node.Id(); present {
				if *idPtr == "sig2_and_sig3a" {
					require.True(t, sig3sent)
					// we've reached the desired outcome
					break loop1
				}
			}
		default:
			t.Logf("%#v", trace)

		}

	}
}
