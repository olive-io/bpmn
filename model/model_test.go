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

package model_test

import (
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/olive-io/bpmn/schema"
	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/model"
	"github.com/olive-io/bpmn/process"
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
	LoadTestFile("testdata/sample.bpmn", &sampleDoc)
}

func TestFindProcess(t *testing.T) {
	model := model.New(&sampleDoc)
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
