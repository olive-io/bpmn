package main

import (
	"context"
	"embed"
	"encoding/xml"
	"log"

	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/process"
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

	for _, processElement := range *definitions.Processes() {
		proc := process.New(&processElement, &definitions)
		if instance, err := proc.Instantiate(); err == nil {
			traces := instance.Tracer.Subscribe()
			err = instance.StartAll(context.Background())
			if err != nil {
				log.Fatalf("failed to run the instance: %s", err)
			}
		loop:
			for {
				trace := tracing.Unwrap(<-traces)
				switch trace := trace.(type) {
				case flow.Trace:
					if id, present := trace.Source.Id(); present {
						if *id == "task" {
							// success!
							break loop
						}

					}
				case activity.ActiveTaskTrace:
					log.Printf("%#v\n", trace)
					trace.Execute()
				case tracing.ErrorTrace:
					log.Fatalf("%#v", trace)
				default:
					log.Printf("%#v\n", trace)
				}
			}
			instance.Tracer.Unsubscribe(traces)
		} else {
			log.Fatalf("failed to instantiate the process: %s", err)
		}
	}
}
