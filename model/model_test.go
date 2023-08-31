package model

import (
	"testing"

	"github.com/olive-io/bpmn/process"
	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/test"
	"github.com/stretchr/testify/assert"
)

func exactId(s string) func(p *process.Process) bool {
	return func(p *process.Process) bool {
		if id, present := p.Element.Id(); present {
			return *id == s
		} else {
			return false
		}
	}
}

var sampleDoc schema.Definitions

func init() {
	test.LoadTestFile("sample/model/sample.bpmn", &sampleDoc)
}

func TestFindProcess(t *testing.T) {
	model := New(&sampleDoc)
	if proc, found := model.FindProcessBy(exactId("sample")); found {
		if id, present := proc.Element.Id(); present {
			assert.Equal(t, *id, "sample")
		} else {
			t.Fatalf("found a process but it has no FlowNodeId")
		}
	} else {
		t.Fatalf("can't find process `sample`")
	}

	if _, found := model.FindProcessBy(exactId("none")); found {
		t.Fatalf("found a process by a non-existent id")
	}
}
