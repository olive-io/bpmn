Introduce ｜ [中文](https://github.com/olive-io/bpmn/tree/main/README_ZH.md)

**Lightweight Business Process Model and Notation (BPMN) 2.0 workflow engine implemented purely in Go**

[![Go Reference](https://pkg.go.dev/badge/github.com/olive-io/bpmn.svg)](https://pkg.go.dev/github.com/olive-io/bpmn)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache-blue.svg)](LICENSE.md)
[![Build Status](https://github.com/olive-io/bpmn/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/olive-io/bpmn/actions/workflows/main.yml?query=branch%3Amain)
[![Last Commit](https://img.shields.io/github/last-commit/olive-io/bpmn)](https://github.com/olive-io/bpmn/commits/main)

---

## Introduction

`github.com/olive-io/bpmn` is a lightweight **Business Process Model and Notation (BPMN) 2.0** workflow engine implemented in Go, designed for modeling and executing business processes embedded in Go applications.  

It supports core BPMN 2.0 elements including Activities (User Task, Service Task, Script Task), Events (Start Event, End Event, Intermediate Catch Event), Gateways (Exclusive Gateway, Inclusive Gateway, Parallel Gateway, Event-based Gateway), Sub-processes, Sequence Flows, and extensible process attributes.

---

## Key Features

- **Standard BPMN 2.0 Compliance** – Build process models using standardized Business Process Model and Notation elements, fully compliant with OMG specification.
- **Lightweight & Embeddable** – Minimal dependencies; easily embedded into business systems with zero external service requirements.
- **Complete Activity Support** – User Tasks, Service Tasks, Script Tasks, Manual Tasks, Business Rule Tasks, and custom Activity implementations.
- **Comprehensive Flow Control** – Sub-processes, Parallel Gateways, Exclusive Gateways, Inclusive Gateways, and Event-based Gateways for complex decision logic.
- **Event-Driven Processing** – Start Events, End Events, Intermediate Catch Events, Timer Events, and Boundary Events.
- **Extensible Process Attributes** – Custom properties and data objects to meet diverse business process requirements.
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
			select {
			case trace := <-traces:
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
			default:
                
			}
		}
	}()
	ins.WaitUntilComplete(ctx)

	pros := ins.Locator().CloneVariables()
	log.Printf("%#v", pros)
}

```

### More Examples
- [Basic Task](https://github.com/olive-io/bpmn/tree/main/examples/basic): Simple Activity execution example
- [User Task](https://github.com/olive-io/bpmn/tree/main/examples/user_task): Implementing User Task Activities
- [Gateways](https://github.com/olive-io/bpmn/tree/main/examples/gateway): Gateway flow control examples
- [Gateways-expr](https://github.com/olive-io/bpmn/tree/main/examples/gateway_expr): Inclusive Gateway with expression evaluation
- [Custom Properties](https://github.com/olive-io/bpmn/tree/main/examples/properties): Process-specific data attributes and custom parameters
- [Catch Event](https://github.com/olive-io/bpmn/tree/main/examples/catch_event): Intermediate Catch Event and Boundary Event examples
- [Throw Event and Collaboration](https://github.com/olive-io/bpmn/tree/main/examples/collaboration): Intermediate Catch Event and Collaboration examples
- [Sub-process](https://github.com/olive-io/bpmn/tree/main/examples/subprocess): Embedded Sub-process execution
- [Multiprocess](https://github.com/olive-io/bpmn/tree/main/examples/multiprocess): Executes multiprocess in an definitions

## License

This project is licensed under the Apache-2.0 License. Commercial and derivative works are welcome.