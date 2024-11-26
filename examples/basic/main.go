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
	"io"
	"log"

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
	instances []*bpmn.Instance
}

func NewWorkflow(reader io.Reader) (*Workflow, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var definitions schema.Definitions
	if err = xml.Unmarshal(data, &definitions); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	tracer := tracing.NewTracer(ctx)

	processes := make([]*bpmn.Process, 0)
	instances := make([]*bpmn.Instance, 0)

	for _, element := range *definitions.Processes() {
		able, ok := element.IsExecutable()
		if !ok {
			continue
		}

		if !able {
			continue
		}

		pr := bpmn.NewProcess(&element, &definitions, bpmn.WithContext(ctx), bpmn.WithTracer(tracer))
		processes = append(processes, pr)

		instance, err := pr.Instantiate()
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
		tracer:      tracer,
		processes:   processes,
		instances:   instances,
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
	_ = ctx
	traces := w.tracer.Subscribe()

	defer w.tracer.Unsubscribe(traces)
	defer w.cancel()

	for _, instance := range w.instances {
		if err := instance.StartAll(); err != nil {
			return err
		}

	LOOP:
		for {
			trace := tracing.Unwrap(<-traces)

			handle(trace)
			if _, ok := trace.(bpmn.CeaseFlowTrace); ok {
				break LOOP
			}
		}

		for !instance.WaitUntilComplete(ctx) {
		}

		//time.Sleep(time.Second * 3)
	}

	return nil
}

func main() {
	var err error
	src, err := fs.Open("task.bpmn")
	if err != nil {
		log.Fatalf("Can't read bpmn: %v", err)
	}

	wf, err := NewWorkflow(src)
	if err != nil {
		log.Fatalf("Can't create workflow: %v", err)
	}

	err = wf.Run(func(trace tracing.ITrace) {
		switch tr := trace.(type) {
		case *bpmn.TaskTrace:
			log.Printf("%#v\n", trace)
			tr.Do()
		case bpmn.ErrorTrace:
			log.Fatalf("%#v", trace)
		default:
			log.Printf("%#v\n", trace)
		}
	})
	if err != nil {
		log.Fatalf("Can't run workflow: %v", err)
	}
}
