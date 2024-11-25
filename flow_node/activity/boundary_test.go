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

package activity_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/flow_node/event/catch"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/tracing"
)

//go:embed testdata
var testdata embed.FS

func LoadTestFile(filename string, definitions any) {
	var err error
	src, err := testdata.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't read file %s: %v", filename, err)
	}
	err = xml.Unmarshal(src, definitions)
	if err != nil {
		log.Fatalf("XML unmarshalling error in %s: %v", filename, err)
	}
}

var testDoc schema.Definitions

func init() {
	LoadTestFile("testdata/boundary_event.bpmn", &testDoc)
}

func TestInterruptingEvent(t *testing.T) {
	testBoundaryEvent(t, "sig1listener", func(visited map[string]bool) {
		assert.False(t, visited["uninterrupted"])
		assert.True(t, visited["interrupted"])
		assert.True(t, visited["end"])
	}, event.NewSignalEvent("sig1"))
}

func TestNonInterruptingEvent(t *testing.T) {
	testBoundaryEvent(t, "sig2listener", func(visited map[string]bool) {
		assert.False(t, visited["interrupted"])
		assert.True(t, visited["uninterrupted"])
		assert.True(t, visited["end"])
	}, event.NewSignalEvent("sig2"))
}

func testBoundaryEvent(t *testing.T, boundary string, test func(visited map[string]bool), events ...event.IEvent) {
	processElement := (*testDoc.Processes())[0]
	proc := process.New(&processElement, &testDoc)
	ready := make(chan bool, 1)

	// explicit tracer
	tracer := tracing.NewTracer(context.Background())
	// this gives us some room when instance starts up
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 32))

	if inst, err := proc.Instantiate(instance.WithTracer(tracer)); err == nil {
		err = inst.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}

		listening := false
		activeBoundary := false

		for {
			if listening && activeBoundary {
				break
			}
			trace := tracing.Unwrap(<-traces)
			t.Logf("%#v", trace)
			switch trace := trace.(type) {
			case catch.ActiveListeningTrace:
				// Ensure the boundary event listener is actually listening
				if id, present := trace.Node.Id(); present {
					if *id == boundary {
						// it is indeed listening
						listening = true
					}
				}
			case activity.ActiveBoundaryTrace:
				if id, present := trace.Node.Id(); present && trace.Start {
					if *id == "task" {
						// task has reached its active boundary
						activeBoundary = true
					}
				}
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				t.Logf("%#v", trace)
			}
		}

		for _, evt := range events {
			_, err = inst.ConsumeEvent(evt)
			assert.Nil(t, err)
		}
		visited := make(map[string]bool)
	loop1:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case flow.VisitTrace:
				if id, present := trace.Node.Id(); present {
					if *id == "uninterrupted" {
						// we're here to we can release the task
						ready <- true
					}
					visited[*id] = true
					if *id == "end" {
						break loop1
					}
				}
			case *activity.Trace:
				if id, present := trace.GetActivity().Element().Id(); present {
					if *id == "task" {
						//v.Execute()
						go func() {
							select {
							case <-ready:
								trace.Do()
							}
						}()
					} else {
						trace.Do()
					}
				} else {
				}
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				t.Logf("%#v", trace)
			}
		}
		inst.Tracer.Unsubscribe(traces)

		test(visited)

	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
