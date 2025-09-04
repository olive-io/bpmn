Introduce ｜ [中文](https://github.com/olive-io/bpmn/tree/main/README_ZH.md)
# github.com/olive-io/bpmn

**Lightweight BPMN 2.0 workflow engine implemented purely in Go**

[![Go Reference](https://pkg.go.dev/badge/github.com/olive-io/bpmn.svg)](https://pkg.go.dev/github.com/olive-io/bpmn)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache-blue.svg)](LICENSE.md)
[![Build Status](https://github.com/olive-io/bpmn/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/olive-io/bpmn/actions/workflows/main.yml?query=branch%3Amain)
[![Last Commit](https://img.shields.io/github/last-commit/olive-io/bpmn)](https://github.com/olive-io/bpmn/commits/main)

---

## Introduction

`github.com/olive-io/bpmn` is a lightweight **BPMN 2.0** workflow engine implemented in Go, designed to simplify modeling and execution of business processes embedded in Go applications.  
It supports core BPMN elements, including task types (User Task, Service Task, Script Task), events (Start / End / Catch), gateways (Exclusive, Inclusive, Parallel, Event-based), sub-processes, flow control, and customizable attributes.

---

## Key Features

- **Native BPMN 2.0 Support** – Build workflows directly with standard BPMN elements, fully compliant with the specification.
- **Lightweight & Simple** – Minimal dependencies; easily embedded into business systems with no extra services required.
- **Multiple Task Types** – User tasks, script tasks, service tasks, and custom tasks.
- **Rich Flow Control** – Sub-processes, parallel, exclusive, and event-based decision logic.
- **Customizable Extensions** – Extend process attributes to meet diverse business scenarios.
- **Comprehensive Test Coverage** – Examples and modules come with unit tests to ensure reliable execution.

---

## Getting Started

### Installation

```bash
go get -u github.com/olive-io/bpmn/schema
go get -u github.com/olive-io/bpmn/v2
```

### Quick Start Example
```go
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

	var err error
	data, err := os.ReadFile("task.bpmn")
	if err != nil {
		log.Fatalf("Can't read bpmn: %v", err)
	}
	definitions, err := schema.Parse(data)
	if err != nil {
		log.Fatalf("XML unmarshalling error: %v", err)
	}

	engine := bpmn.NewEngine()
	options := []bpmn.Option{
		bpmn.WithVariables(map[string]any{
			"c": map[string]string{"name": "cc"},
		}),
		bpmn.WithDataObjects(map[string]any{
			"a": struct{}{},
		}),
	}
	ctx := context.Background()
	ins, err := engine.NewProcess(&definitions, options...)
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
			}

			trace = tracing.Unwrap(trace)
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
			case bpmn.TaskTrace:
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
}

```

### More Examples
- [Basic Task](https://github.com/olive-io/bpmn/tree/main/examples/basic): Simplest task example
- [User Task](https://github.com/olive-io/bpmn/tree/main/examples/user_task): Executing a user task
- [Gateways](https://github.com/olive-io/bpmn/tree/main/examples/gateway): Gateway execution examples
- [Gateways-expr](https://github.com/olive-io/bpmn/tree/main/examples/gateway_expr): Inclusive Gateway with Expr
- [Custom Properties](https://github.com/olive-io/bpmn/tree/main/examples/properties): Support for custom task parameters
- [Sub-process](https://github.com/olive-io/bpmn/tree/main/examples/subprocess): Executing a sub-process

## License

This project is licensed under the Apache-2.0 License. Commercial and derivative works are welcome.