/*
   Copyright 2025 The bpmn Authors

   This program is offered under a commercial and under the AGPL license.
   For AGPL licensing, see below.

   AGPL licensing:
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"embed"
	"encoding/xml"
	"log"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/event"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

//go:embed boundary_event.bpmn
var fs embed.FS

func main() {
	var definitions schema.Definitions

	var err error
	src, err := fs.ReadFile("boundary_event.bpmn")
	if err != nil {
		log.Fatalf("Can't read bpmn: %v", err)
	}
	err = xml.Unmarshal(src, &definitions)
	if err != nil {
		log.Fatalf("XML unmarshalling error: %v", err)
	}

	engine := bpmn.NewEngine()
	options := []bpmn.Option{}
	ctx := context.Background()
	ins, err := engine.NewProcess(&definitions, options...)
	if err != nil {
		log.Fatalf("failed to instantiate the process: %s", err)
		return
	}
	traces := ins.Tracer().Subscribe()
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
			case bpmn.TaskTrace:
				act := trace.GetActivity()
				ele := act.Element()
				name, _ := ele.Name()
				log.Printf("Do Task [%s]", *name)
				trace.Do()
			case bpmn.ActiveListeningTrace:
				for _, eventDefinition := range trace.Node.SignalEventDefinitionField {
					ref, ok := eventDefinition.SignalRef()
					if ok {
						ins.ConsumeEvent(event.NewSignalEvent(string(*ref)))
					}
				}
				for _, eventDefinition := range trace.Node.MessageEventDefinitionField {
					ref, ok := eventDefinition.MessageRef()
					if ok {
						ins.ConsumeEvent(event.NewMessageEvent(string(*ref), (*string)(eventDefinition.OperationRefField)))
					}
				}

			case bpmn.ErrorTrace:
				log.Fatalf("%#v", trace)
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
