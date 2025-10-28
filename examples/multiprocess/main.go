/*
Copyright 2025 The bpmn Authors

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

//go:embed sample.bpmn
var fs embed.FS

func main() {
	var definitions schema.Definitions

	var err error
	src, err := fs.ReadFile("sample.bpmn")
	if err != nil {
		log.Fatalf("Can't read bpmn: %v", err)
	}
	err = xml.Unmarshal(src, &definitions)
	if err != nil {
		log.Fatalf("XML unmarshalling error: %v", err)
	}

	engine := bpmn.NewEngine()
	var options []bpmn.Option
	ctx := context.Background()
	ins, err := engine.NewProcessSet(&definitions, options...)
	if err != nil {
		log.Fatalf("failed to instantiate the process: %s", err)
		return
	}
	traces := ins.Tracer().Subscribe()
	defer ins.Tracer().Unsubscribe(traces)
	err = ins.StartAll(ctx)
	if err != nil {
		log.Fatalf("failed to run the instance: %s", err)
	}
	go func() {
		for {
			var trace tracing.ITrace
			select {
			case trace = <-traces:
			default:
				continue
			}

			trace = tracing.Unwrap(trace)
			switch trace := trace.(type) {
			case bpmn.TaskTrace:
				act := trace.GetActivity()
				ele := act.Element()
				name, _ := ele.Name()
				log.Printf("Do Task [%s]", *name)
				trace.Do()
			case bpmn.ErrorTrace:
				log.Fatalf("%#v", trace)
			case bpmn.CeaseProcessSetTrace:
				return
			default:
				var name string
				v, ok := trace.Unpack().(schema.FlowNodeInterface)
				if ok {
					value, found := v.Name()
					if found {
						name = *value
					}
				}
				log.Printf("%#v: %s", trace, name)
			}
		}
	}()
	ins.WaitUntilComplete(ctx)

	pros := ins.Locator().CloneVariables()
	log.Printf("%#v", pros)
}
