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

	processElement := (*definitions.Processes())[0]
	proc := bpmn.NewProcess(&processElement, &definitions)
	options := []bpmn.Option{
		bpmn.WithVariables(map[string]any{}),
		bpmn.WithDataObjects(map[string]any{}),
	}
	ctx := context.Background()
	if ins, err := proc.Instantiate(options...); err == nil {
		traces := ins.Tracer().Subscribe()
		err = ins.StartAll()
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
				case *bpmn.TaskTrace:
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

		pros := ins.Locator().CloneVariables()
		log.Printf("%#v", pros)
		ins.Tracer().Unsubscribe(traces)
	} else {
		log.Fatalf("failed to instantiate the process: %s", err)
	}
}
