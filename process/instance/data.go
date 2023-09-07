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

package instance

import (
	"sync"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/schema"
)

type DataObjectContainer struct {
	data.DefaultItemAwareLocator

	mu                         sync.RWMutex
	dataObjectsByName          map[string]data.IItemAware
	dataObjects                map[schema.Id]data.IItemAware
	dataObjectReferencesByName map[string]data.IItemAware
	dataObjectReferences       map[schema.Id]data.IItemAware
	propertiesByName           map[string]data.IItemAware
	properties                 map[schema.Id]data.IItemAware
}

func NewDataObjectContainer() *DataObjectContainer {
	return &DataObjectContainer{
		dataObjectsByName:          map[string]data.IItemAware{},
		dataObjects:                map[schema.Id]data.IItemAware{},
		dataObjectReferencesByName: map[string]data.IItemAware{},
		dataObjectReferences:       map[schema.Id]data.IItemAware{},
		propertiesByName:           map[string]data.IItemAware{},
		properties:                 map[schema.Id]data.IItemAware{},
	}
}

func (do *DataObjectContainer) FindItemAwareById(id schema.IdRef) (itemAware data.IItemAware, found bool) {
	do.mu.RLock()
	defer do.mu.RUnlock()
	for k := range do.dataObjects {
		if k == id {
			found = true
			itemAware = do.dataObjects[k]
			goto ready
		}
	}
	for k := range do.dataObjectReferences {
		if k == id {
			found = true
			itemAware = do.dataObjectReferences[k]
			goto ready
		}
	}
	for k := range do.properties {
		if k == id {
			found = true
			itemAware = do.properties[k]
			goto ready
		}
	}
ready:
	return
}

func (do *DataObjectContainer) FindItemAwareByName(name string) (itemAware data.IItemAware, found bool) {
	do.mu.RLock()
	defer do.mu.RUnlock()
	for k := range do.dataObjectsByName {
		if k == name {
			found = true
			itemAware = do.dataObjectsByName[k]
			goto ready
		}
	}
	for k := range do.dataObjectReferencesByName {
		if k == name {
			found = true
			itemAware = do.dataObjectReferencesByName[k]
			goto ready
		}
	}
	for k := range do.propertiesByName {
		if k == name {
			found = true
			itemAware = do.propertiesByName[k]
			goto ready
		}
	}
ready:
	return
}

func (do *DataObjectContainer) PutItemAwareByName(name string, itemAware data.IItemAware) {
	do.mu.Lock()
	defer do.mu.Unlock()
	do.dataObjectsByName[name] = itemAware
}

func (do *DataObjectContainer) Clone() map[string]any {
	out := make(map[string]any)
	for name, item := range do.dataObjectsByName {
		value := item.Get()
		if value != nil {
			out[name] = value
		}
	}
	for name, item := range do.propertiesByName {
		value := item.Get()
		if value != nil {
			out[name] = value
		}
	}
	return out
}

type HeaderContainer struct {
	data.DefaultItemAwareLocator
	items map[string]data.IItemAware
}

func NewHeaderContainer() *HeaderContainer {
	return &HeaderContainer{
		items: map[string]data.IItemAware{},
	}
}

func (h *HeaderContainer) FindItemAwareById(id schema.IdRef) (data.IItemAware, bool) {
	return nil, false
}

func (h *HeaderContainer) FindItemAwareByName(name string) (data.IItemAware, bool) {
	item, ok := h.items[name]
	if !ok {
		return nil, false
	}

	return item, true
}

func (h *HeaderContainer) Clone() map[string]any {
	out := make(map[string]any)
	for name, item := range h.items {
		value := item.Get()
		if value != nil {
			out[name] = value
		}
	}
	return out
}

type PropertyContainer struct {
	mu sync.RWMutex
	data.DefaultItemAwareLocator
	items map[string]data.IItemAware
}

func NewPropertyContainer() *PropertyContainer {
	return &PropertyContainer{
		items: map[string]data.IItemAware{},
	}
}

func (p *PropertyContainer) FindItemAwareById(id schema.IdRef) (data.IItemAware, bool) {
	return nil, false
}

func (p *PropertyContainer) FindItemAwareByName(name string) (data.IItemAware, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	item, ok := p.items[name]
	if !ok {
		return nil, false
	}

	return item, true
}

func (p *PropertyContainer) PutItemAwareByName(name string, itemAware data.IItemAware) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.items[name] = itemAware
}

func (p *PropertyContainer) Clone() map[string]any {
	out := make(map[string]any)
	for name, item := range p.items {
		value := item.Get()
		if value != nil {
			out[name] = value
		}
	}
	return out
}

type FlowDataLocator struct {
	lmu       sync.RWMutex
	locators  map[string]data.IItemAwareLocator
	vmu       sync.RWMutex
	variables map[string]data.IItem
}

func NewFlowDataLocator() *FlowDataLocator {
	f := &FlowDataLocator{
		locators:  map[string]data.IItemAwareLocator{},
		variables: map[string]data.IItem{},
	}
	return f
}

func (f *FlowDataLocator) FindIItemAwareLocator(name string) (locator data.IItemAwareLocator, found bool) {
	f.lmu.RLock()
	defer f.lmu.RUnlock()
	locator, found = f.locators[name]
	return locator, found
}

func (f *FlowDataLocator) PutIItemAwareLocator(name string, locator data.IItemAwareLocator) {
	f.lmu.Lock()
	defer f.lmu.Unlock()
	f.locators[name] = locator
}

func (f *FlowDataLocator) CloneItems(name string) map[string]any {
	out := make(map[string]any)

	f.vmu.RLock()
	locator, ok := f.locators[name]
	if !ok {
		f.vmu.RUnlock()
		return out
	}
	f.vmu.RUnlock()

	return locator.Clone()
}

func (f *FlowDataLocator) GetVariable(name string) (value any, found bool) {
	f.vmu.RLock()
	defer f.vmu.RUnlock()
	value, found = f.variables[name]
	return
}

func (f *FlowDataLocator) SetVariable(name string, value any) {
	f.vmu.Lock()
	defer f.vmu.Unlock()
	f.variables[name] = value
}

func (f *FlowDataLocator) CloneVariables() map[string]any {
	f.vmu.RLock()
	defer f.vmu.RUnlock()
	out := make(map[string]any)
	for key, value := range f.variables {
		out[key] = value
	}
	return out
}
