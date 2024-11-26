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

	"github.com/stretchr/testify/require"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/model"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
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
			_, ok := trace.(bpmn.FlowTrace)
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
		case bpmn.VisitTrace:
			if idPtr, present := trace.Node.Id(); present {
				if *idPtr == "sig1a" {
					// we've reached the desired outcome
					break loop1
				}
			}
		default:
			//t.Logf("%#v", trace)
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
			_, ok := trace.(bpmn.FlowTrace)
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
		case bpmn.VisitTrace:
			if idPtr, present := trace.Node.Id(); present {
				if *idPtr == "sig2_sig3a" {
					// we've reached the desired outcome
					break loop1
				}
			}
		default:
			//t.Logf("%#v", trace)
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
			_, ok := trace.(bpmn.FlowTrace)
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
		case bpmn.VisitTrace:
			if idPtr, present := trace.Node.Id(); present {
				if *idPtr == "sig2_and_sig3a" {
					require.True(t, sig3sent)
					// we've reached the desired outcome
					break loop1
				}
			}
		default:
			//t.Logf("%#v", trace)

		}

	}
}
