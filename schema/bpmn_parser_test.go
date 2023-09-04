// Copyright 2023 Lack (xingyys@gmail.com).
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

func TestParseSample(t *testing.T) {
	var sampleDoc Definitions
	var err error
	LoadTestFile("testdata/sample.bpmn", &sampleDoc)
	processes := sampleDoc.Processes()
	assert.Equal(t, 1, len(*processes))

	out, err := xml.MarshalIndent(&sampleDoc, "", " ")
	if !assert.NoError(t, err) {
		return
	}

	t.Log(string(out))
}

func TestParseSampleNs(t *testing.T) {
	var sampleDoc Definitions
	LoadTestFile("testdata/sample_ns.bpmn", &sampleDoc)
	processes := sampleDoc.Processes()
	assert.Equal(t, 1, len(*processes))
}
