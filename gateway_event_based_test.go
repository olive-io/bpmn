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

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func TestEventBasedGateway(t *testing.T) {
	var testDoc schema.Definitions
	LoadTestFile("testdata/event_based_gateway.bpmn", &testDoc)

	testEventBasedGateway(t, func(reached map[string]int) {
		assert.Equal(t, 1, reached["task1"])
		assert.Equal(t, 1, reached["end"])
	}, event.NewSignalEvent("Sig1"))
	testEventBasedGateway(t, func(reached map[string]int) {
		assert.Equal(t, 1, reached["task2"])
		assert.Equal(t, 1, reached["end"])
	}, event.NewMessageEvent("Msg1", nil))
}

func testEventBasedGateway(t *testing.T, test func(map[string]int), events ...event.IEvent) {
	var testDoc schema.Definitions
	LoadTestFile("testdata/event_based_gateway.bpmn", &testDoc)

	engine := bpmn.NewEngine()
	ctx := context.Background()
	proc, err := engine.NewProcess(&testDoc)
	if err == nil {
		traces := proc.Tracer().Subscribe()
		err := proc.StartAll(ctx)
		if err != nil {
			t.Errorf("failed to run the instance: %s", err)
			return
		}

		resultChan := make(chan string)

		ctx, cancel := context.WithCancel(context.Background())

		go func(ctx context.Context) {
			for {
				select {
				case trace := <-traces:
					trace = tracing.Unwrap(trace)
					switch trace := trace.(type) {
					case bpmn.ActiveListeningTrace:
						if id, present := trace.Node.Id(); present {
							resultChan <- *id
						}
					case bpmn.ErrorTrace:
						t.Errorf("%#v", trace)
						return
					default:
						//t.Logf("%#v", trace)
					}
				case <-ctx.Done():
					return
				}
			}
		}(ctx)

		// Wait until both events are ready to listen
		assert.Regexp(t, "(signalEvent|messageEvent)", <-resultChan)
		assert.Regexp(t, "(signalEvent|messageEvent)", <-resultChan)

		cancel()

		ch := make(chan map[string]int)
		go func() {
			reached := make(map[string]int)
			for {
				var trace tracing.ITrace

				select {
				case trace = <-traces:
					trace = tracing.Unwrap(trace)
				default:
					continue
				}
				switch trace := trace.(type) {
				case bpmn.VisitTrace:
					if id, present := trace.Node.Id(); present {
						if counter, ok := reached[*id]; ok {
							reached[*id] = counter + 1
						} else {
							reached[*id] = 1
						}
					} else {
						t.Errorf("can't find element with FlowNodeId %#v", id)
						ch <- reached
						return
					}
				case bpmn.TaskTrace:
					trace.Do()
				case bpmn.CeaseFlowTrace:
					ch <- reached
					return
				case bpmn.ErrorTrace:
					t.Errorf("%#v", trace)
					ch <- reached
					return
				default:
					//t.Logf("%#v", trace)
				}
			}
		}()

		for _, evt := range events {
			_, err := proc.ConsumeEvent(evt)
			if err != nil {
				t.Error(err)
				return
			}
		}

		test(<-ch)

		proc.Tracer().Unsubscribe(traces)
	} else {
		t.Errorf("failed to instantiate the process: %s", err)
		return
	}
}
