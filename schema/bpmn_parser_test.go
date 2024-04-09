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

package schema

import (
	"embed"
	"encoding/xml"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata
var testdata embed.FS

func LoadTestFile(filename string) *Definitions {
	var err error
	src, err := testdata.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't read file %s: %v", filename, err)
	}
	definitions, err := Parse(src)
	if err != nil {
		log.Fatalf("Can't parse file %s: %v", filename, err)
	}
	return definitions
}

func TestParseSample(t *testing.T) {
	sampleDoc := LoadTestFile("testdata/sample.bpmn")
	processes := sampleDoc.Processes()
	assert.Equal(t, 1, len(*processes))

	_, err := xml.MarshalIndent(&sampleDoc, "", " ")
	if !assert.NoError(t, err) {
		return
	}

	element, found := sampleDoc.FindBy(ExactId("left"))
	if !assert.True(t, found) {
		return
	}

	id, _ := element.(BaseElementInterface).Id()
	name, _ := element.(FlowElementInterface).Name()
	t.Log(*id, *name)

	element, found = sampleDoc.FindBy(ExactId("right"))
	if !assert.True(t, found) {
		return
	}

	extension, ok := element.(FlowNodeInterface).ExtensionElements()
	if !assert.True(t, ok) {
		return
	}

	t.Log(extension.TaskHeaderField.Header[0])
}

func TestParseSampleNs(t *testing.T) {
	sampleDoc := LoadTestFile("testdata/sample_ns.bpmn")
	processes := sampleDoc.Processes()
	assert.Equal(t, 1, len(*processes))
}
