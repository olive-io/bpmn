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
