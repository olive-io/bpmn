BPMN 2.0 Schema Package
=======================

A lightweight, pure Go implementation of a BPMN 2.0 parser
----------------------------------------------------------

This package provides complete serialization and deserialization functionality for the BPMN 2.0 XML format, supporting all standard BPMN 2.0 elements and extensions.

Features
--------
- 🚀 Pure Go implementation with no external dependencies
- 📋 Complete support for the BPMN 2.0 specification
- 🔄 XML serialization/deserialization
- 🔍 Element lookup and filtering
- 🎯 Type-safe element access
- 🛠️ Support for custom extensions (Olive extension)

Quick Start
-----------

### Installation

go get -u github.com/olive-io/bpmn/schema

### Basic Usage

#### Parsing a BPMN File

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/olive-io/bpmn/schema"
)

func main() {
	// Read the BPMN file
	data, err := os.ReadFile("process.bpmn")
	if err != nil {
		log.Fatalf("Unable to read BPMN file: %v", err)
	}

	// Parse the BPMN definitions
	definitions, err := schema.Parse(data)
	if err != nil {
		log.Fatalf("XML parsing error: %v", err)
	}

	fmt.Printf("Parsed BPMN definitions: %s\n", definitions.Name())
}

```

#### Building Definitions and Auto Layout (BPMN DI)

```go
package main

import (
	"encoding/xml"
	"fmt"

	"github.com/olive-io/bpmn/schema"
)

func main() {
	// Build process graph with fluent builders.
	pb := schema.NewProcessBuilder()
	task1 := schema.Task{}
	task2 := schema.ScriptTask{}
	pb.AddActivity(&task1).AddActivity(&task2)

	db := schema.NewDefinitionsBuilder()
	db.AddProcess(*pb.Out())

	// Generate BPMN DI (BPMNShape/BPMNEdge/Bounds/waypoints).
	// The default layout keeps node centers and sequence-flow lines aligned.
	db.AutoLayout(schema.DefaultAutoLayoutConfig())

	// Optional: customize start point and spacing.
	db.AutoLayout(&schema.AutoLayoutConfig{
		StartX:     96,
		StartY:     96,
		ColumnGap:  220,
		RowGap:     140,
		ProcessGap: 180,
	})

	definitions := db.Out()
	data, _ := xml.MarshalIndent(definitions, "", "  ")
	fmt.Println(string(data))
}
```

Notes:

- Auto layout writes `definitions.DiagramField` with `BPMNDiagram/BPMNPlane`.
- For each flow node, a `BPMNShape` with `dc:Bounds` is generated.
- For each sequence flow, a `BPMNEdge` with aligned `di:waypoint` list is generated.
- If multiple processes are added, participants are generated under one collaboration.

#### Element Lookup

```go
// Find element by ID
element, found := definitions.FindBy(schema.ExactId("task-1"))
if found {
    if task, ok := element.(*schema.ServiceTask); ok {
        fmt.Printf("Found service task: %s\n", task.Name())
    }
}

// Find element by type
startEvent, found := definitions.FindBy(schema.ElementType(&schema.StartEvent{}))
if found {
    fmt.Printf("Found start event: %s\n", startEvent.(*schema.StartEvent).Name())
}

// Find element by interface type
var flowNodeInterface schema.FlowNodeInterface
flowNode, found := definitions.FindBy(schema.ElementInterface(&flowNodeInterface))
if found {
    fmt.Printf("Found flow node: %s\n", flowNode.(schema.FlowNodeInterface).Name())
}

```

#### Combined Queries

```go
// Combine query conditions using And and Or operators
task, found := definitions.FindBy(
        schema.ElementType(&schema.Task{}).And(
        schema.ExactId("my-task"),
    ),
)

```

### XML Serialization
```go
// Create new definitions
definitions := &schema.Definitions{
    IdField:              schema.NewStringP("def-1"),
    NameField:            schema.NewStringP("Sample Process"),
    TargetNamespaceField: "http://example.com/bpmn",
}

// Serialize to XML
data, err := xml.Marshal(definitions)
if err != nil {
    log.Fatalf("XML serialization error: %v", err)
}

fmt.Println(string(data))
```


Core Types
----------

### Base Types

- Element - The base interface for all BPMN elements
- BaseElementInterface - The interface for elements that have an ID
- FlowNodeInterface - The interface for flow nodes
- ElementPredicate - The function type for an element predicate

### Main Elements

#### Definitions and Processes
- Definitions - The root element of BPMN definitions
- Process - A business process definition
- Collaboration - A collaboration definition

#### Flow Nodes
- StartEvent - Start event
- EndEvent - End event
- IntermediateThrowEvent - Intermediate throw event
- IntermediateCatchEvent - Intermediate catch event
- Task - Base task
- ServiceTask - Service task
- UserTask - User task
- ScriptTask - Script task
- BusinessRuleTask - Business rule task
- CallActivity - Call activity
- SubProcess - Sub-process

#### Gateways
- ExclusiveGateway - Exclusive gateway
- ParallelGateway - Parallel gateway
- InclusiveGateway - Inclusive gateway
- EventBasedGateway - Event-based gateway

#### Connectors
- SequenceFlow - Sequence flow
- MessageFlow - Message flow
- Association - Association

Utility Functions
-----------------

#### Pointer Constructors

```go
// Pointers for basic types
stringPtr := schema.NewStringP("value")
intPtr := schema.NewIntegerP(42)
boolPtr := schema.NewBoolP(true)
floatPtr := schema.NewFloatP(3.14)

// QName type
qname := schema.NewQName("prefix:localName")

#### Lookup Predicates

// Exact ID match
predicate := schema.ExactId("element-id")

// Element type match
predicate := schema.ElementType(&schema.Task{})

// Interface type match
var flowNode schema.FlowNodeInterface
predicate := schema.ElementInterface(&flowNode)
```

Extension Support
-----------------

### Olive Extension

```go
This package supports the custom Olive extension to enhance BPMN functionality:

// Task definition extension
taskDef := &schema.TaskDefinition{
    Type:    "http-service",
    Timeout: "30s",
    Retries: 3,
    Target:  "node01",
}

// Properties extension
properties := &schema.Properties{
    Property: []*schema.Item{
        {
            Name:  "endpoint",
            Value: "/api/users",
            Type:  schema.ItemTypeString,
        },
    },
}

// Extension elements
extensions := &schema.ExtensionElementsType{
    TaskDefinitionField: taskDef,
    PropertiesField:     properties,
}
```

### Extension Element Types

- TaskDefinition - Task definition
- TaskHeader - Task header
- Properties - Properties configuration
- Result - Result configuration
- ExtensionScript - Script extension
- ExtensionCalledElement - Called element extension
- ExtensionCalledDecision - Called decision extension

Data Type Support
-----------------

### Item Type System
```go

// Supported data types
const (
    ItemTypeObject  = "object"   // JSON object
    ItemTypeArray   = "array"    // JSON array
    ItemTypeInteger = "integer"  // Integer
    ItemTypeString  = "string"   // String
    ItemTypeBoolean = "boolean"  // Boolean
    ItemTypeFloat   = "float"    // Float
)

// Create an Item
itemObject := &schema.Item{  
    Name:  "config",
    Value: `{"timeout": 30, "retries": 3}`,
    Type:  schema.ItemTypeObject,
}

itemArray := &schema.Item{
    Name:  "config",
    Value: `["1", "2", "3"]`,
    Type:  schema.ItemTypeArray,
}


// Get the typed value
value := item.ValueFor() // Returns map[string]any

// Bind to a struct
type Config struct {
    Timeout int `json:"timeout"`
    Retries int `json:"retries"`
}
var config Config
err := itemObject.ValueTo(&config)

```
Advanced Features
-----------------

### Expression Support

```go
// Supports formal expressions and regular expressions with DataObject
expressionDataObject := &schema.AnExpression{
	Expression: &schema.FormalExpression{
        TextPayloadField: schema.NewStringP("getDataObject('dataObject1').number > 10"),
    },
}

// Supports formal expressions and regular expressions with Headers
expressionHeader := &schema.AnExpression{
    Expression: &schema.FormalExpression{
        TextPayloadField: schema.NewStringP("getHeader('age') > 10"),
    },
}

// Supports formal expressions and regular expressions with Properties
expressionProperties := &schema.AnExpression{
    Expression: &schema.FormalExpression{
        TextPayloadField: schema.NewStringP("getProp('aa') > 10"),
    },
}
```

### Process Instantiation

```go
// Get the flow nodes that can instantiate a process
instantiatingNodes := process.InstantiatingFlowNodes()
for _, node := range instantiatingNodes {
    fmt.Printf("Instantiating node: %s\n", node.Name())
}
```


Best Practices
--------------

### 1. Error Handling

```go
data, _ := os.ReadFile("task.bpmn")
definitions, err := schema.Parse(data)
if err != nil {
    // Handle parsing errors
    log.Printf("BPMN parsing failed: %v", err)
    return
}
```

### 2. Type Assertion

```go
element, found := definitions.FindBy(schema.ExactId("task-1"))
if found {
    if serviceTask, ok := element.(*schema.ServiceTask); ok {
        // Safely use serviceTask
    }
}
```

### 3. Memory Management

```go
// Use pointer constructors to avoid manually creating pointers
task := &schema.Task{
    IdField:   schema.NewStringP("task-1"),
    NameField: schema.NewStringP("Sample Task"),
}

```
License
-------

Copyright 2023 The bpmn Authors

Licensed under the Apache License, Version 2.0
