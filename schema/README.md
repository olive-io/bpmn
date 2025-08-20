## Introduce

**Lightweight BPMN 2.0 parser implemented purely in Go**

## Getting Started

### Installation

```bash
go get -u github.com/olive-io/bpmn/schema
```

### Quick Start Example

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/olive-io/bpmn/schema"
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

	element, found := definitions.FindBy(schema.ExactId("left"))
	if found {
		task := element.(*schema.ServiceTask)
		fmt.Println(task.Id())
	}
}
```