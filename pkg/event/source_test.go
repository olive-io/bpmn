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

package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoidSource_RegisterEventConsumer(t *testing.T) {
	source := VoidSource{}
	consumer := &mockConsumer{result: Consumed}

	err := source.RegisterEventConsumer(consumer)

	assert.NoError(t, err)
	// VoidSource should not actually register the consumer,
	// it just returns nil error
}

func TestVoidSource_InterfaceImplementation(t *testing.T) {
	source := VoidSource{}
	
	// Verify it implements ISource interface
	var iSource ISource = source
	assert.NotNil(t, iSource)
	
	// Test the interface method
	consumer := &mockConsumer{result: Consumed}
	err := iSource.RegisterEventConsumer(consumer)
	assert.NoError(t, err)
}

func TestISource_Interface(t *testing.T) {
	// Test that various implementations satisfy the ISource interface
	
	t.Run("VoidSource implements ISource", func(t *testing.T) {
		var source ISource = VoidSource{}
		assert.Implements(t, (*ISource)(nil), source)
	})
	
	t.Run("FanOut implements ISource", func(t *testing.T) {
		var source ISource = NewFanOut()
		assert.Implements(t, (*ISource)(nil), source)
	})
}