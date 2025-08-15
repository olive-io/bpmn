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

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/tracing"

	"github.com/stretchr/testify/assert"
)

func TestSignalEvent(t *testing.T) {
	testEvent(t, "testdata/intermediate_catch_event.bpmn", "signalCatch", nil, false, event.NewSignalEvent("global_sig1"))
}

func TestMessageEvent(t *testing.T) {
	testEvent(t, "testdata/intermediate_catch_event.bpmn", "messageCatch", nil, false, event.NewMessageEvent("msg", nil))
}

func TestMultipleEvent(t *testing.T) {
	// either
	testEvent(t, "testdata/intermediate_catch_event_multiple.bpmn", "multipleCatch", nil, false, event.NewMessageEvent("msg", nil))
	// or
	testEvent(t, "testdata/intermediate_catch_event_multiple.bpmn", "multipleCatch", nil, false, event.NewSignalEvent("global_sig1"))
}

func TestMultipleParallelEvent(t *testing.T) {
	// both
	testEvent(t, "testdata/intermediate_catch_event_multiple_parallel.bpmn", "multipleParallelCatch",
		nil, false, event.NewMessageEvent("msg", nil), event.NewSignalEvent("global_sig1"))
	// either
	testEvent(t, "testdata/intermediate_catch_event_multiple_parallel.bpmn", "multipleParallelCatch", nil, true, event.NewMessageEvent("msg", nil))
	testEvent(t, "testdata/intermediate_catch_event_multiple_parallel.bpmn", "multipleParallelCatch", nil, true, event.NewSignalEvent("global_sig1"))
}

type eventInstance struct {
	id string
}

func (e eventInstance) EventDefinition() schema.EventDefinitionInterface {
	definition := schema.DefaultEventDefinition()
	return &definition
}

type eventDefinitionInstanceBuilder struct{}

func (e eventDefinitionInstanceBuilder) NewEventDefinitionInstance(def schema.EventDefinitionInterface) (event.IDefinitionInstance, error) {
	switch d := def.(type) {
	case *schema.TimerEventDefinition:
		id, _ := d.Id()
		return eventInstance{id: *id}, nil
	case *schema.ConditionalEventDefinition:
		id, _ := d.Id()
		return eventInstance{id: *id}, nil
	default:
		return event.WrapEventDefinition(d), nil
	}
}

func TestTimerEvent(t *testing.T) {
	i := eventInstance{id: "timer_1"}
	b := eventDefinitionInstanceBuilder{}
	testEvent(t, "testdata/intermediate_catch_event.bpmn", "timerCatch", &b, false, event.MakeTimerEvent(i))
}

func TestConditionalEvent(t *testing.T) {
	i := eventInstance{id: "conditional_1"}
	b := eventDefinitionInstanceBuilder{}
	testEvent(t, "testdata/intermediate_catch_event.bpmn", "conditionalCatch", &b, false, event.MakeTimerEvent(i))
}

func testEvent(t *testing.T, filename string, nodeId string, eventDefinitionInstanceBuilder event.IDefinitionInstanceBuilder, eventObservationOnly bool, events ...event.IEvent) {
	var testDoc schema.Definitions
	LoadTestFile(filename, &testDoc)
	engine := bpmn.NewEngine()

	tracer := tracing.NewTracer(context.Background())
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 64))

	inst, err := engine.NewProcess(&testDoc, bpmn.WithProcessEventDefinitionInstanceBuilder(eventDefinitionInstanceBuilder), bpmn.WithTracer(tracer))
	assert.Nil(t, err)

	err = inst.StartAll()
	if err != nil {
		t.Fatalf("failed to run the instance: %s", err)
	}
	resultChan := make(chan bool)
	go func() {
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.ActiveListeningTrace:
				if id, present := trace.Node.Id(); present {
					if *id == nodeId {
						// listening
						resultChan <- true
						return
					}

				}
			case bpmn.ErrorTrace:
				t.Errorf("%#v", trace)
				resultChan <- false
				return
			default:
				//t.Logf("%#v", trace)
			}
		}
	}()

	assert.True(t, <-resultChan)

	go func() {
		defer inst.Tracer().Unsubscribe(traces)
		eventsToObserve := events
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case bpmn.EventObservedTrace:
				if eventObservationOnly {
					for i := range eventsToObserve {
						if eventsToObserve[i] == trace.Event {
							eventsToObserve[i] = eventsToObserve[len(eventsToObserve)-1]
							eventsToObserve = eventsToObserve[:len(eventsToObserve)-1]
							break
						}
					}
					if len(eventsToObserve) == 0 {
						resultChan <- true
						return
					}
				}
			case bpmn.FlowTrace:
				if id, present := trace.Source.Id(); present {
					if *id == nodeId {
						// success!
						resultChan <- true
						return
					}

				}
			case bpmn.ErrorTrace:
				t.Errorf("%#v", trace)
				resultChan <- false
				return
			default:
				//t.Logf("%#v", trace)
			}
		}
	}()

	for _, evt := range events {
		_, err = inst.ConsumeEvent(evt)
		assert.Nil(t, err)
	}

	assert.True(t, <-resultChan)

}
