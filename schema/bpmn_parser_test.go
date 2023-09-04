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
