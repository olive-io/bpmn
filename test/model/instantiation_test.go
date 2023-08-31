package model

import (
	"context"
	"testing"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/model"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/test"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/require"
)

var testStartEventInstantiation schema.Definitions

func init() {
	test.LoadTestFile("sample/model/instantiate_start_event.bpmn", &testStartEventInstantiation)
}

func TestStartEventInstantiation(t *testing.T) {
	ctx := context.Background()
	tracer := tracing.NewTracer(ctx)
	traces := tracer.SubscribeChannel(make(chan tracing.Trace, 128))
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
	traces := tracer.SubscribeChannel(make(chan tracing.Trace, 128))
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
	traces := tracer.SubscribeChannel(make(chan tracing.Trace, 128))
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
