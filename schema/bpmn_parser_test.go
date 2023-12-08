// Copyright 2023 The olive Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	out, err := xml.MarshalIndent(&sampleDoc, "", " ")
	if !assert.NoError(t, err) {
		return
	}

	t.Log(string(out))
}

func TestParseSampleNs(t *testing.T) {
	sampleDoc := LoadTestFile("testdata/sample_ns.bpmn")
	processes := sampleDoc.Processes()
	assert.Equal(t, 1, len(*processes))
}
