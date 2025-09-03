/*
Copyright 2023 The bpmn Authors

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

package id

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarningTrace_Unpack(t *testing.T) {
	t.Run("string warning", func(t *testing.T) {
		warning := "test warning message"
		trace := WarningTrace{Warning: warning}
		
		result := trace.Unpack()
		assert.Equal(t, warning, result)
	})

	t.Run("struct warning", func(t *testing.T) {
		type TestWarning struct {
			Code    int
			Message string
		}
		
		warning := TestWarning{Code: 500, Message: "internal error"}
		trace := WarningTrace{Warning: warning}
		
		result := trace.Unpack()
		assert.Equal(t, warning, result)
		
		unpacked, ok := result.(TestWarning)
		assert.True(t, ok)
		assert.Equal(t, 500, unpacked.Code)
		assert.Equal(t, "internal error", unpacked.Message)
	})

	t.Run("nil warning", func(t *testing.T) {
		trace := WarningTrace{Warning: nil}
		
		result := trace.Unpack()
		assert.Nil(t, result)
	})

	t.Run("int warning", func(t *testing.T) {
		warning := 42
		trace := WarningTrace{Warning: warning}
		
		result := trace.Unpack()
		assert.Equal(t, warning, result)
	})

	t.Run("map warning", func(t *testing.T) {
		warning := map[string]interface{}{
			"error":   "overflow",
			"context": "sequence generation",
		}
		trace := WarningTrace{Warning: warning}
		
		result := trace.Unpack()
		assert.Equal(t, warning, result)
		
		unpacked, ok := result.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "overflow", unpacked["error"])
		assert.Equal(t, "sequence generation", unpacked["context"])
	})
}