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

	"github.com/olive-io/bpmn/schema"

	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/process"
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

	for _, processElement := range *definitions.Processes() {
		proc := process.New(&processElement, &definitions)
		if instance, err := proc.Instantiate(); err == nil {
			traces := instance.Tracer.Subscribe()
			ctx, cancel := context.WithCancel(context.Background())
			err = instance.StartAll(ctx)
			if err != nil {
				cancel()
				log.Fatalf("failed to run the instance: %s", err)
			}
			go func() {
			LOOP:
				for {
					trace := tracing.Unwrap(<-traces)
					switch trace := trace.(type) {
					case *activity.Trace:
						log.Printf("%#v\n", trace)
						trace.Do()
					case tracing.ErrorTrace:
						log.Fatalf("%#v", trace)
					case flow.CeaseFlowTrace:
						break LOOP
					default:
						log.Printf("%#v\n", trace)
					}
				}
			}()
			instance.WaitUntilComplete(ctx)
			instance.Tracer.Unsubscribe(traces)
			cancel()
		} else {
			log.Fatalf("failed to instantiate the process: %s", err)
		}
	}
}
