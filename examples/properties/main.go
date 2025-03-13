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

	xml.Marshal(definitions)

	processElement := (*definitions.Processes())[0]
	proc := bpmn.NewProcess(&processElement, definitions)
	options := []bpmn.Option{
		bpmn.WithVariables(map[string]any{
			"c": map[string]string{"name": "cc"},
		}),
		bpmn.WithDataObjects(map[string]any{
			"a": struct{}{},
		}),
	}
	ctx := context.Background()
	ins, err := proc.Instantiate(options...)
	if err != nil {
		log.Fatalf("failed to instantiate the process: %s", err)
		return
	}
	traces := ins.Tracer().Subscribe()
	err = ins.StartAll()
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
			case bpmn.FlowTrace:
			case *bpmn.TaskTrace:
				trace.Do(bpmn.DoWithResults(
					map[string]any{
						"c": map[string]string{"name": "cc1"},
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
	ins.WaitUntilComplete(ctx)

	pros := ins.Locator().CloneVariables()
	log.Printf("%#v", pros)
	ins.Tracer().Unsubscribe(traces)
}
