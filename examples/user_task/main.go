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

package main

import (
	"context"
	"embed"
	"encoding/xml"
	"log"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

//go:embed user_task.bpmn
var fs embed.FS

func main() {
	var definitions schema.Definitions

	var err error
	src, err := fs.ReadFile("user_task.bpmn")
	if err != nil {
		log.Fatalf("Can't read bpmn: %v", err)
	}
	err = xml.Unmarshal(src, &definitions)
	if err != nil {
		log.Fatalf("XML unmarshalling error: %v", err)
	}

	cache := map[string]struct{}{}
	users := map[string]string{}

	engine := bpmn.NewEngine()

	options := []bpmn.Option{
		bpmn.WithVariables(map[string]any{}),
		bpmn.WithDataObjects(map[string]any{}),
	}

	proc, err := engine.NewProcess(&definitions, options...)
	if err != nil {
		log.Fatalf("Can't create process: %v", err)
	}
	ctx := context.Background()
	traces := proc.Tracer().Subscribe()
	defer proc.Tracer().Unsubscribe(traces)
	err = proc.StartAll()
	if err != nil {
		log.Fatalf("failed to run the instance: %s", err)
	}
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		for {
			var trace tracing.ITrace
			select {
			case trace = <-traces:
			}

			trace = tracing.Unwrap(trace)
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
			case bpmn.TaskTrace:
				id, _ := trace.GetActivity().Element().Id()
				if _, ok := cache[*id]; ok {
					// already executed, skip it
					break
				}

				cache[*id] = struct{}{}

				uid := trace.GetProperties()["uid"]
				users[uid.(string)] = "waiting"

				//TODO: waiting for client requesting

				// executes user task
				trace.Do()
				log.Printf("%#v", trace)
			case bpmn.ErrorTrace:
				log.Fatalf("%#v", trace)
				return
			case bpmn.CeaseFlowTrace:
				return
			default:
				log.Printf("%#v", trace)
			}
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}

	proc.WaitUntilComplete(ctx)

	//pros := ins.Locator().CloneVariables()
	log.Printf("%#v", users)
}
