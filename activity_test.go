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
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
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

func TestInterruptingEvent(t *testing.T) {
	var testDoc schema.Definitions
	LoadTestFile("testdata/boundary_event.bpmn", &testDoc)

	testBoundaryEvent(t, "sig1listener", func(visited map[string]bool) {
		assert.False(t, visited["uninterrupted"])
		assert.True(t, visited["interrupted"])
		assert.True(t, visited["end"])
		t.Log(visited)
	}, event.NewSignalEvent("sig1"))
}

func TestNonInterruptingEvent(t *testing.T) {
	var testDoc schema.Definitions
	LoadTestFile("testdata/boundary_event.bpmn", &testDoc)

	testBoundaryEvent(t, "sig2listener", func(visited map[string]bool) {
		assert.False(t, visited["interrupted"])
		assert.True(t, visited["uninterrupted"])
		assert.True(t, visited["end"])
		t.Log(visited)
	}, event.NewSignalEvent("sig2"))
}

func testBoundaryEvent(t *testing.T, boundary string, test func(visited map[string]bool), events ...event.IEvent) {
	var testDoc schema.Definitions
	LoadTestFile("testdata/boundary_event.bpmn", &testDoc)

	engine := bpmn.NewEngine()
	ready := make(chan bool, 1)

	ctx := context.Background()
	// explicit tracer
	tracer := tracing.NewTracer(ctx)
	// this gives us some room when instance starts up
	traces := tracer.SubscribeChannel(make(chan tracing.ITrace, 32))
	defer tracer.Unsubscribe(traces)

	if inst, err := engine.NewProcess(&testDoc, bpmn.WithTracer(tracer)); err == nil {
		err = inst.StartAll(ctx)
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
			//t.Logf("%#v", trace)
			switch trace := trace.(type) {
			case bpmn.ActiveListeningTrace:
				// Ensure the boundary event listener is actually listening
				if id, present := trace.Node.Id(); present {
					if *id == boundary {
						// it is indeed listening
						listening = true
					}
				}
			case bpmn.ActiveBoundaryTrace:
				if id, present := trace.Node.Id(); present && trace.Start {
					if *id == "task" {
						// task has reached its active boundary
						activeBoundary = true
					}
				}
			case bpmn.ErrorTrace:
				t.Fatalf("%#v", trace)
			default:
				//t.Logf("%#v", trace)
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
			case bpmn.VisitTrace:
				if id, present := trace.Node.Id(); present {
					if *id == "uninterrupted" {
						// we're here to we can release the task
						ready <- true
					}
					if *id == "interrupted" {
						ready <- true
					}
					visited[*id] = true
					if *id == "end" {
						break loop1
					}
				}
			case bpmn.TaskTrace:
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
			case bpmn.ErrorTrace:
				// t.Fatalf("%#v", trace)
			default:
				//t.Logf("%#v", trace)
			}
		}

		test(visited)

	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
