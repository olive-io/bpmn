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
	"io"
	"log"
	"time"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

//go:embed task.bpmn
var fs embed.FS

type Workflow struct {
	ctx    context.Context
	cancel context.CancelFunc

	definitions *schema.Definitions
	tracer      tracing.ITracer
	//traces      chan tracing.ITrace
	processes []*bpmn.Process
}

func NewWorkflow(ctx context.Context, reader io.Reader, opts ...bpmn.Option) (*Workflow, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var definitions schema.Definitions
	if err = xml.Unmarshal(data, &definitions); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	tracer := tracing.NewTracer(ctx)

	instances := make([]*bpmn.Process, 0)

	opts = append(opts, bpmn.WithContext(ctx), bpmn.WithTracer(tracer))

	for _, element := range *definitions.Processes() {
		able, ok := element.IsExecutable()
		if !ok || !able {
			continue
		}

		instance, err := bpmn.NewProcess(&element, &definitions, opts...)
		if err != nil {
			cancel()
			return nil, err
		}
		instances = append(instances, instance)
	}

	workflow := &Workflow{
		ctx:         ctx,
		cancel:      cancel,
		definitions: &definitions,
		processes:   instances,
		tracer:      tracer,
	}

	return workflow, nil
}

//func (w *Workflow) Trace() chan tracing.ITrace {
//	if w.traces != nil {
//		w.traces = w.tracer.Subscribe()
//	}
//
//	return w.traces
//}

type Handle func(trace tracing.ITrace)

func (w *Workflow) Run(handle Handle) error {
	ctx := w.ctx
	traces := w.tracer.Subscribe()

	defer w.tracer.Unsubscribe(traces)
	defer w.cancel()

	for _, instance := range w.processes {
		if err := instance.StartAll(); err != nil {
			return err
		}

	LOOP:
		for {
			wrapped, ok := <-traces
			if !ok {
				break LOOP
			}

			trace := tracing.Unwrap(wrapped)

			handle(trace)
			if _, ok := trace.(bpmn.CeaseFlowTrace); ok {
				break LOOP
			}
		}

		instance.WaitUntilComplete(ctx)
	}

	return nil
}

func main() {
	var err error
	src, err := fs.Open("task.bpmn")
	if err != nil {
		log.Fatalf("Can't read bpmn: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	wf, err := NewWorkflow(ctx, src)
	if err != nil {
		log.Fatalf("Can't create workflow: %v", err)
	}

	err = wf.Run(func(trace tracing.ITrace) {
		switch tr := trace.(type) {
		case bpmn.TaskTrace:
			log.Printf("%#v\n", trace)
			tr.Do()
		case bpmn.ErrorTrace:
			log.Fatalf("%#v", trace)
		default:
			if tr == nil {
				log.Fatalf("empty trace: %v", tr)
			}
			log.Printf("%#v\n", trace)
		}
	})
	if err != nil {
		log.Fatalf("Can't run workflow: %v", err)
	}
}
