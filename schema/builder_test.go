/*
Copyright 2025 The bpmn Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package schema

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandBytes(t *testing.T) {
	t.Run("RandBytes", func(t *testing.T) {
		out := RandBytes(2)
		assert.Equal(t, len(out), 2)
		out1 := RandBytes(2)
		assert.NotEqual(t, out, out1)
		t.Log("out:", string(out))

		out2 := RandBytes(10)
		assert.NotEqual(t, out, out2)
		t.Log("out:", string(out2))
	})
}

func TestDefinitionsWithProcess(t *testing.T) {
	db := NewDefinitionsBuilder()
	pb := NewProcessBuilder()
	task := Task{}
	pb.AddActivity(&task)
	script := ScriptTask{}
	pb.AddActivity(&script)

	db.AddProcess(*pb.Out())

	def := db.Out()

	assert.Equal(t, len(def.ProcessField), 1)

	data, _ := xml.MarshalIndent(def, "", " ")
	t.Log(string(data))
}
