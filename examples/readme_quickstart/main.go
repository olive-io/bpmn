package main

import (
	"context"
	"log"
	"os"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func main() {
	data, err := os.ReadFile("task.bpmn")
	if err != nil {
		log.Fatalf("read bpmn file: %v", err)
	}

	definitions, err := schema.Parse(data)
	if err != nil {
		log.Fatalf("parse bpmn xml: %v", err)
	}

	engine := bpmn.NewEngine()
	ctx := context.Background()

	proc, err := engine.NewProcess(definitions,
		bpmn.WithVariables(map[string]any{"customer": "alice"}),
	)
	if err != nil {
		log.Fatalf("create process: %v", err)
	}

	traces := proc.Tracer().Subscribe()
	defer proc.Tracer().Unsubscribe(traces)

	if err := proc.StartAll(ctx); err != nil {
		log.Fatalf("start process: %v", err)
	}

	go func() {
		for trace := range traces {
			trace = tracing.Unwrap(trace)
			switch t := trace.(type) {
			case bpmn.TaskTrace:
				t.Do(bpmn.DoWithResults(map[string]any{"approved": true}))
			case bpmn.ErrorTrace:
				log.Printf("process error: %v", t.Error)
			}
		}
	}()

	if ok := proc.WaitUntilComplete(ctx); !ok {
		log.Printf("process cancelled")
	}

	log.Printf("variables: %#v", proc.Locator().CloneVariables())
}
