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
	"github.com/olive-io/bpmn/schema"
)

type mockEventDefinition struct {
	eventType         string
	id                *schema.Id
	documentations    *[]schema.Documentation
	extensionElements *schema.ExtensionElements
}

func (m *mockEventDefinition) Clone() interface{} {
	return &mockEventDefinition{
		eventType:         m.eventType,
		id:                m.id,
		documentations:    m.documentations,
		extensionElements: m.extensionElements,
	}
}

func (m *mockEventDefinition) FindBy(predicate schema.ElementPredicate) (result schema.Element, found bool) {
	if predicate(m) {
		return m, true
	}
	return nil, false
}

func (m *mockEventDefinition) Id() (result *schema.Id, present bool) {
	return m.id, m.id != nil
}

func (m *mockEventDefinition) Documentations() (result *[]schema.Documentation) {
	return m.documentations
}

func (m *mockEventDefinition) ExtensionElements() (result *schema.ExtensionElements, present bool) {
	return m.extensionElements, m.extensionElements != nil
}

func (m *mockEventDefinition) SetId(value *schema.Id) {
	m.id = value
}

func (m *mockEventDefinition) SetDocumentations(value []schema.Documentation) {
	m.documentations = &value
}

func (m *mockEventDefinition) SetExtensionElements(value *schema.ExtensionElements) {
	m.extensionElements = value
}

type mockDefinitionInstance struct {
	definition schema.EventDefinitionInterface
}

func (m *mockDefinitionInstance) EventDefinition() schema.EventDefinitionInterface {
	return m.definition
}

type mockConsumer struct {
	result   ConsumptionResult
	err      error
	consumed []IEvent
}

func (m *mockConsumer) ConsumeEvent(ev IEvent) (ConsumptionResult, error) {
	m.consumed = append(m.consumed, ev)
	return m.result, m.err
}