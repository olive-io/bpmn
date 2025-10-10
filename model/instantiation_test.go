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

package model_test

//
//var testStartEventInstantiation schema.Definitions
//
//func init() {
//	LoadTestFile("testdata/instantiate_start_event.bpmn", &testStartEventInstantiation)
//}
//
//func TestStartEventInstantiation(t *testing.T) {
//	ctx := context.Background()
//	tracer := tracing.NewTracer(ctx)
//	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
//	m, err := model.New(&testStartEventInstantiation, model.WithContext(ctx), model.WithTracer(tracer))
//	require.Nil(t, err)
//	err = m.Run(ctx)
//	require.Nil(t, err)
//loop:
//	for {
//		select {
//		case trace := <-traces:
//			trace = tracing.Unwrap(trace)
//			_, ok := trace.(bpmn.FlowTrace)
//			// Should not flow
//			require.False(t, ok)
//		default:
//			break loop
//		}
//	}
//	// Test simple event instantiation
//	_, err = m.ConsumeEvent(event.NewSignalEvent("sig1"))
//	require.Nil(t, err)
//loop1:
//	for {
//		var trace tracing.ITrace

//select {
//case trace = <-traces:
//	trace = tracing.Unwrap(trace)
//default:
//	continue
//}
//		switch trace := trace.(type) {
//		case bpmn.VisitTrace:
//			if idPtr, present := trace.Node.Id(); present {
//				if *idPtr == "sig1a" {
//					// we've reached the desired outcome
//					break loop1
//				}
//			}
//		default:
//			//t.Logf("%#v", trace)
//		}
//	}
//}
//
//func TestMultipleStartEventInstantiation(t *testing.T) {
//	ctx := context.Background()
//	tracer := tracing.NewTracer(ctx)
//	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
//	m, err := model.New(&testStartEventInstantiation, model.WithContext(ctx), model.WithTracer(tracer))
//	require.Nil(t, err)
//	err = m.Run(ctx)
//	require.Nil(t, err)
//loop:
//	for {
//		select {
//		case trace := <-traces:
//			trace = tracing.Unwrap(trace)
//			_, ok := trace.(bpmn.FlowTrace)
//			// Should not flow
//			require.False(t, ok)
//		default:
//			break loop
//		}
//	}
//	// Test multiple event instantiation
//	_, err = m.ConsumeEvent(event.NewSignalEvent("sig2"))
//	require.Nil(t, err)
//loop1:
//	for {
//		var trace tracing.ITrace
//
//			select {
//			case trace = <-traces:
//				trace = tracing.Unwrap(trace)
//			default:
//				continue
//			}
//		switch trace := trace.(type) {
//		case bpmn.VisitTrace:
//			if idPtr, present := trace.Node.Id(); present {
//				if *idPtr == "sig2_sig3a" {
//					// we've reached the desired outcome
//					break loop1
//				}
//			}
//		default:
//			//t.Logf("%#v", trace)
//		}
//	}
//}
//
//func TestParallelMultipleStartEventInstantiation(t *testing.T) {
//	ctx := context.Background()
//	tracer := tracing.NewTracer(ctx)
//	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 128))
//	m, err := model.New(&testStartEventInstantiation, model.WithContext(ctx), model.WithTracer(tracer))
//	require.Nil(t, err)
//	err = m.Run(ctx)
//	require.Nil(t, err)
//loop:
//	for {
//		select {
//		case trace := <-traces:
//			trace = tracing.Unwrap(trace)
//			_, ok := trace.(bpmn.FlowTrace)
//			// Should not flow
//			require.False(t, ok)
//		default:
//			break loop
//		}
//	}
//	// Test parallel multiple event instantiation
//	signalEvent := event.NewSignalEvent("sig2")
//	_, err = m.ConsumeEvent(signalEvent)
//	require.Nil(t, err)
//	sig3sent := false
//loop1:
//	for {
//		var trace tracing.ITrace
//
//			select {
//			case trace = <-traces:
//				trace = tracing.Unwrap(trace)
//			default:
//				continue
//			}
//		switch trace := trace.(type) {
//		case model.EventInstantiationAttemptedTrace:
//			if !sig3sent && signalEvent == trace.Event {
//				if idPtr, present := trace.Node.Id(); present && *idPtr == "ParallelMultipleStartEvent" {
//					_, err = m.ConsumeEvent(event.NewSignalEvent("sig3"))
//					require.Nil(t, err)
//					sig3sent = true
//				}
//			}
//		case bpmn.VisitTrace:
//			if idPtr, present := trace.Node.Id(); present {
//				if *idPtr == "sig2_and_sig3a" {
//					require.True(t, sig3sent)
//					// we've reached the desired outcome
//					break loop1
//				}
//			}
//		default:
//			//t.Logf("%#v", trace)
//
//		}
//
//	}
//}
