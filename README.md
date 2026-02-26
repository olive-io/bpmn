Introduce | [中文](https://github.com/olive-io/bpmn/tree/main/README_ZH.md)

**Lightweight BPMN 2.0 workflow engine written in Go**

[![Go Reference](https://pkg.go.dev/badge/github.com/olive-io/bpmn.svg)](https://pkg.go.dev/github.com/olive-io/bpmn)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache-blue.svg)](LICENSE.md)
[![Build Status](https://github.com/olive-io/bpmn/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/olive-io/bpmn/actions/workflows/main.yml?query=branch%3Amain)
[![Last Commit](https://img.shields.io/github/last-commit/olive-io/bpmn)](https://github.com/olive-io/bpmn/commits/main)

---

## What Is This?

`github.com/olive-io/bpmn` is an embeddable BPMN 2.0 runtime for Go applications.

Use it when you want to:

- Execute BPMN process definitions inside your Go service.
- Handle user/service/script tasks with your own business logic.
- Track process runtime via trace streams for debugging and observability.

It includes two modules:

- Runtime module: `github.com/olive-io/bpmn/v2`
- BPMN schema module: `github.com/olive-io/bpmn/schema`

---

## Installation

```bash
go get -u github.com/olive-io/bpmn/schema
go get -u github.com/olive-io/bpmn/v2
```

---

## Quick Start

This example loads a BPMN file, starts a process, handles task traces, and waits for completion.

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
				// Complete user/service/script task with results.
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
```

---

## Core Concepts

- `Engine`: creates process instances from BPMN definitions.
- `Process`: executes one executable process.
- `ProcessSet`: executes multiple executable processes in one definition.
- `Tracer`: runtime event stream; all traces are accessible via subscribe/unsubscribe.
- `TaskTrace`: callback point to provide task result, data objects, headers, and errors.

---

## Single Process vs Process Set

Use `NewProcess` when your BPMN definitions contain exactly one executable process.

Use `NewProcessSet` when definitions contain multiple executable processes.

```go
proc, err := engine.NewProcess(definitions)
// Returns an error if multiple executable processes are present.

set, err := engine.NewProcessSet(&definitions)
// Runs all executable processes in the same definitions.
```

---

## Observability and Runtime Errors

Subscribe tracer events to inspect runtime behavior:

- `FlowTrace`: sequence flow transitions.
- `TaskTrace`: task is waiting for external decision/result.
- `ErrorTrace`: runtime error happened.
- `CeaseFlowTrace` / `CeaseProcessSetTrace`: process or process set completed.

ID generator fallback behavior:

- Runtime attempts to initialize the default SNO-based ID generator.
- If initialization fails, runtime falls back to a local generator.
- A warning trace is emitted to make this visible in observability pipelines.

---

## BPMN Coverage

This engine supports core BPMN elements used by most service workflows, including:

- Tasks: User, Service, Script, Manual, Business Rule, Call Activity.
- Events: Start, End, Intermediate Catch/Throw, Timer-related flows.
- Gateways: Exclusive, Inclusive, Parallel, Event-based.
- Sub-processes and sequence flows.

For concrete usage patterns, see examples and tests in this repository.

---

## Development Commands

From repository root:

```bash
# all runtime tests
go test -v ./...

# schema module tests
go test -v ./schema

# static analysis
go vet ./...

# run one test by name
go test -v ./... -run 'TestUserTask'

# run one test in a package, bypass cache
go test -v ./pkg/data -run '^TestContainer$' -count=1
```

---

## Examples

- [quickstart](https://github.com/olive-io/bpmn/tree/main/examples/readme_quickstart)
- [basic](https://github.com/olive-io/bpmn/tree/main/examples/basic)
- [user_task](https://github.com/olive-io/bpmn/tree/main/examples/user_task)
- [gateway](https://github.com/olive-io/bpmn/tree/main/examples/gateway)
- [gateway_expr](https://github.com/olive-io/bpmn/tree/main/examples/gateway_expr)
- [properties](https://github.com/olive-io/bpmn/tree/main/examples/properties)
- [catch_event](https://github.com/olive-io/bpmn/tree/main/examples/catch_event)
- [collaboration](https://github.com/olive-io/bpmn/tree/main/examples/collaboration)
- [subprocess](https://github.com/olive-io/bpmn/tree/main/examples/subprocess)
- [multiprocess](https://github.com/olive-io/bpmn/tree/main/examples/multiprocess)

---

## Contributing

PRs are welcome. For contributor/agent workflow guidance, see `AGENTS.md`.

---

## License

Apache-2.0
