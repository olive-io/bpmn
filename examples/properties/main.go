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
	"log"
	"time"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

//go:embed task.bpmn
var fs embed.FS

func main() {

	var err error
	src, err := fs.ReadFile("task.bpmn")
	if err != nil {
		log.Fatalf("Can't read bpmn: %v", err)
	}
	definitions, err := schema.Parse(src)
	if err != nil {
		log.Fatalf("XML unmarshalling error: %v", err)
	}

	options := []bpmn.Option{
		bpmn.WithVariables(map[string]any{
			"d": "hello",
			"c": map[string]string{"name": "cc"},
		}),
		bpmn.WithDataObjects(map[string]any{
			"a": struct{}{},
		}),
	}
	engine := bpmn.NewEngine()
	ins, err := engine.NewProcess(definitions, options...)
	if err != nil {
		log.Fatalf("Can't create process: %v", err)
	}

	ctx := context.Background()
	traces := ins.Tracer().Subscribe()
	defer ins.Tracer().Unsubscribe(traces)
	err = ins.StartAll(ctx)
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
				trace = tracing.Unwrap(trace)
			default:
				continue
			}

			switch trace := trace.(type) {
			case bpmn.FlowTrace:
			case bpmn.TaskTrace:
				time.Sleep(1 * time.Second)
				trace.Do(bpmn.DoWithResults(
					map[string]any{
						//"c": map[string]string{"name": "cc1"},
						"a": 2,
					}),
				)
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
	completed := ins.WaitUntilComplete(ctx)
	select {
	case <-done:
	}

	target := map[string]any{}
	err = ins.Locator().ApplyTo(&target)
	if err != nil {
		log.Fatalf("Can't locate target: %v", err)
	}
	log.Printf("%#v: completed: %v", target["c"], completed)
}
