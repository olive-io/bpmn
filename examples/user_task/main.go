package main

import (
	"context"
	"embed"
	"encoding/xml"
	"log"

	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity/user"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/tracing"
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
	proc := process.New(&processElement, &definitions)
	options := []instance.Option{
		instance.WithVariables(map[string]any{}),
		instance.WithDataObjects(map[string]any{}),
	}
	ctx := context.Background()
	if ins, err := proc.Instantiate(options...); err == nil {
		traces := ins.Tracer.Subscribe()
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
				}

				trace = tracing.Unwrap(trace)
				switch trace := trace.(type) {
				case flow.Trace:
				case *user.ActiveTrace:
					id, _ := trace.Activity.Element().Id()
					if _, ok := cache[*id]; ok {
						// already executed, skip it
						break
					}

					cache[*id] = struct{}{}

					uid := trace.Properties["uid"]
					users[uid.(string)] = "waiting"

					//TODO: waiting for client requesting

					// executes user task
					trace.Do()
					log.Printf("%#v", trace)
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

		select {
		case <-done:
		case <-ctx.Done():
		}

		pros := ins.Locator.CloneVariables()
		log.Printf("%#v", pros)
		ins.Tracer.Unsubscribe(traces)
	} else {
		log.Fatalf("failed to instantiate the process: %s", err)
	}
}
