/*
   Copyright 2023 The bpmn Authors

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

package script_test

import (
	"context"
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/olive-io/bpmn/schema"
	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/flow"
	"github.com/olive-io/bpmn/flow_node/activity"
	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/process/instance"
	"github.com/olive-io/bpmn/tracing"
)

//go:embed testdata
var testdata embed.FS

func LoadTestFile(filename string, definitions any) {
	var err error
	src, err := testdata.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't read file %s: %v", filename, err)
	}
	err = xml.Unmarshal(src, definitions)
	if err != nil {
		log.Fatalf("XML unmarshalling error in %s: %v", filename, err)
	}
}

var testTask schema.Definitions

func init() {
	LoadTestFile("testdata/script_task.bpmn", &testTask)
}

func TestScriptTask(t *testing.T) {
	processElement := (*testTask.Processes())[0]
	proc := process.New(&processElement, &testTask)
	option := instance.WithVariables(map[string]any{
		"c": map[string]string{"name": "cc"},
	})
	if instance, err := proc.Instantiate(option); err == nil {
		traces := instance.Tracer.Subscribe()
		err = instance.StartAll(context.Background())
		if err != nil {
			t.Fatalf("failed to run the instance: %s", err)
		}
	loop:
		for {
			trace := tracing.Unwrap(<-traces)
			switch trace := trace.(type) {
			case flow.Trace:
				if id, present := trace.Source.Id(); present {
					if *id == "task" {
						// success!
						//break loop
					}
				}
			case *activity.Trace:
				trace.Do()
				t.Logf("%#v", trace)
			case tracing.ErrorTrace:
				t.Fatalf("%#v", trace)
			case flow.CeaseFlowTrace:
				break loop
			default:
				t.Logf("%#v", trace)
			}
		}
		instance.Tracer.Unsubscribe(traces)
		vv, _ := instance.Locator.GetVariable("sum")
		assert.Equal(t, vv, 4)
	} else {
		t.Fatalf("failed to instantiate the process: %s", err)
	}
}
