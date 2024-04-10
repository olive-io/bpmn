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
