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

package event_based_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/olive-io/bpmn/event"
	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity"
	ev "github.com/olive-io/bpmn/flow_node/event/catch"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
	"github.com/stretchr/testify/assert"
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
	LoadTestFile("testdata/event_based_gateway.bpmn", &testDoc)
}

func TestEventBasedGateway(t *testing.T) {
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
	processElement := (*testDoc.Processes())[0]
	proc := process.New(&processElement, &testDoc)
	if instance, err := proc.Instantiate(); err == nil {
		traces := instance.Tracer.Subscribe()
		err := instance.StartAll(context.Background())
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
					case ev.ActiveListeningTrace:
						if id, present := trace.Node.Id(); present {
							resultChan <- *id
						}
					case tracing.ErrorTrace:
						t.Errorf("%#v", trace)
						return
					default:
						t.Logf("%#v", trace)
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
				trace := tracing.Unwrap(<-traces)
				switch trace := trace.(type) {
				case flow.VisitTrace:
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
				case *activity.Trace:
					trace.Do()
				case flow.CeaseFlowTrace:
					ch <- reached
					return
				case tracing.ErrorTrace:
					t.Errorf("%#v", trace)
					ch <- reached
					return
				default:
					t.Logf("%#v", trace)
				}
			}
		}()

		for _, evt := range events {
			_, err := instance.ConsumeEvent(evt)
			if err != nil {
				t.Error(err)
				return
			}
		}

		test(<-ch)

		instance.Tracer.Unsubscribe(traces)
	} else {
		t.Errorf("failed to instantiate the process: %s", err)
		return
	}
}
