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

import "github.com/olive-io/bpmn/schema"

// IDefinitionInstance is a unifying interface for representing event definition within
// an execution context (useful for event definitions like timer, condition, etc.)
type IDefinitionInstance interface {
	EventDefinition() schema.EventDefinitionInterface
}

// wrappedDefinitionInstance is a simple wrapper for schema.EventDefinitionInterface
// that adds no extra context
type wrappedDefinitionInstance struct {
	definition schema.EventDefinitionInterface
}

func (d *wrappedDefinitionInstance) EventDefinition() schema.EventDefinitionInterface {
	return d.definition
}

// WrapEventDefinition is a default event instance builder that creates Instance simply by
// enclosing schema.EventDefinitionInterface
func WrapEventDefinition(def schema.EventDefinitionInterface) IDefinitionInstance {
	return &wrappedDefinitionInstance{definition: def}
}

// IDefinitionInstanceBuilder allows supplying custom instance builders that interact with the
// rest of the system and add context for further matching
type IDefinitionInstanceBuilder interface {
	NewEventDefinitionInstance(def schema.EventDefinitionInterface) (definitionInstance IDefinitionInstance, err error)
}

type wrappingDefinitionInstanceBuilder struct{}

var WrappingDefinitionInstanceBuilder = wrappingDefinitionInstanceBuilder{}

func (d wrappingDefinitionInstanceBuilder) NewEventDefinitionInstance(def schema.EventDefinitionInterface) (IDefinitionInstance, error) {
	return WrapEventDefinition(def), nil
}

type fallbackDefinitionInstanceBuilder struct {
	builders []IDefinitionInstanceBuilder
}

func (f *fallbackDefinitionInstanceBuilder) NewEventDefinitionInstance(def schema.EventDefinitionInterface) (definitionInstance IDefinitionInstance, err error) {
	for i := range f.builders {
		definitionInstance, err = f.builders[i].NewEventDefinitionInstance(def)
		if err != nil {
			return
		}
		if definitionInstance != nil {
			return
		}
	}
	return
}

// DefinitionInstanceBuildingChain creates a DefinitionInstanceBuilder that attempts supplied builders
// from left to right, until a builder returns a non-nil DefinitionInstanceBuilder, which is then
// returned from the call to DefinitionInstanceBuildingChain
func DefinitionInstanceBuildingChain(builders ...IDefinitionInstanceBuilder) IDefinitionInstanceBuilder {
	builder := &fallbackDefinitionInstanceBuilder{builders: builders}
	return builder
}
