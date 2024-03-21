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

package main

import (
	"context"
	"embed"
	"encoding/xml"
	"log"

	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
)

//go:embed task.bpmn
var fs embed.FS

func main() {
	var definitions schema.Definitions

	var err error
	src, err := fs.ReadFile("task.bpmn")
	if err != nil {
		log.Fatalf("Can't read bpmn: %v", err)
	}
	err = xml.Unmarshal(src, &definitions)
	if err != nil {
		log.Fatalf("XML unmarshalling error: %v", err)
	}

	processElement := (*definitions.Processes())[0]
	proc := process.New(&processElement, &definitions)
	options := []instance.Option{
		instance.WithVariables(map[string]any{
			"c": map[string]string{"name": "cc"},
		}),
		instance.WithDataObjects(map[string]any{
			"a": struct{}{},
		}),
	}
	ctx := context.Background()
	ins, err := proc.Instantiate(options...)
	if err != nil {
		log.Fatalf("failed to instantiate the process: %s", err)
		return
	}
	traces := ins.Tracer.Subscribe()
	err = ins.StartAll(ctx)
	if err != nil {
		log.Fatalf("failed to run the instance: %s", err)
	}
	go func() {
		for {
			var trace tracing.ITrace
			select {
			case trace = <-traces:
			}

			trace = tracing.Unwrap(trace)
			switch trace := trace.(type) {
			case flow.Trace:
			case *activity.Trace:
				trace.Do(activity.WithProperties(
					map[string]any{
						"c": map[string]string{"name": "cc1"},
						"a": 2,
					}),
				)
			case tracing.ErrorTrace:
				log.Fatalf("%#v", trace)
				return
			case flow.CeaseFlowTrace:
				return
			default:
				log.Printf("%#v", trace)
			}
		}
	}()
	ins.WaitUntilComplete(ctx)

	pros := ins.Locator.CloneVariables()
	log.Printf("%#v", pros)
	ins.Tracer.Unsubscribe(traces)
}
