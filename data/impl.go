/*
   Copyright 2023 The bpmn Authors

   This library is free software; you can redistribute it and/or
   modify it under the terms of the GNU Lesser General Public
   License as published by the Free Software Foundation; either
   version 2.1 of the License, or (at your option) any later version.

   This library is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
   Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public
   License along with this library;
*/

package data

import (
	"fmt"
	"sync"

	json "github.com/bytedance/sonic"

	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/pkg/id"
	"github.com/olive-io/bpmn/schema"
)

const (
	LocatorObject   = "$"
	LocatorHeader   = "#"
	LocatorProperty = "@"
)

type ObjectContainer struct {
	DefaultItemAwareLocator

	mu                         sync.RWMutex
	dataObjectsByName          map[string]IItemAware
	dataObjects                map[schema.Id]IItemAware
	dataObjectReferencesByName map[string]IItemAware
	dataObjectReferences       map[schema.Id]IItemAware
	propertiesByName           map[string]IItemAware
	properties                 map[schema.Id]IItemAware
}

func NewDataObjectContainer() *ObjectContainer {
	return &ObjectContainer{
		dataObjectsByName:          map[string]IItemAware{},
		dataObjects:                map[schema.Id]IItemAware{},
		dataObjectReferencesByName: map[string]IItemAware{},
		dataObjectReferences:       map[schema.Id]IItemAware{},
		propertiesByName:           map[string]IItemAware{},
		properties:                 map[schema.Id]IItemAware{},
	}
}

func (do *ObjectContainer) FindItemAwareById(id schema.IdRef) (itemAware IItemAware, found bool) {
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

func (do *ObjectContainer) FindItemAwareByName(name string) (itemAware IItemAware, found bool) {
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

func (do *ObjectContainer) PutItemAwareById(id schema.IdRef, itemAware IItemAware) {
	do.mu.Lock()
	defer do.mu.Unlock()
	do.dataObjects[id] = itemAware
}

func (do *ObjectContainer) PutItemAwareByName(name string, itemAware IItemAware) {
	do.mu.Lock()
	defer do.mu.Unlock()
	do.dataObjectsByName[name] = itemAware
}

func (do *ObjectContainer) Clone() map[string]any {
	out := make(map[string]any)
	for name, item := range do.dataObjects {
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

func (do *ObjectContainer) CloneFor(other ILocatorCloner) {
	out, ok := other.(*ObjectContainer)
	if !ok {
		return
	}

	out.mu.Lock()
	defer out.mu.Unlock()

	do.mu.RLock()
	defer do.mu.RUnlock()

	cloneItemAwareMap(do.dataObjectsByName, &out.dataObjectsByName)
	cloneItemAwareMap(do.dataObjects, &out.dataObjects)
	cloneItemAwareMap(do.dataObjectReferencesByName, &out.dataObjectReferencesByName)
	cloneItemAwareMap(do.dataObjectReferences, &out.dataObjectReferences)
	cloneItemAwareMap(do.propertiesByName, &out.propertiesByName)
	cloneItemAwareMap(do.properties, &out.properties)
}

type HeaderContainer struct {
	DefaultItemAwareLocator
	items map[string]IItemAware
}

func NewHeaderContainer() *HeaderContainer {
	return &HeaderContainer{
		items: map[string]IItemAware{},
	}
}

func (h *HeaderContainer) FindItemAwareById(id schema.IdRef) (IItemAware, bool) {
	return nil, false
}

func (h *HeaderContainer) FindItemAwareByName(name string) (IItemAware, bool) {
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

func (h *HeaderContainer) CloneFor(other ILocatorCloner) {
	out, ok := other.(*HeaderContainer)
	if !ok {
		return
	}

	cloneItemAwareMap(h.items, &out.items)
}

type PropertyContainer struct {
	mu sync.RWMutex
	DefaultItemAwareLocator
	items map[string]IItemAware
}

func NewPropertyContainer() *PropertyContainer {
	return &PropertyContainer{
		items: map[string]IItemAware{},
	}
}

func (p *PropertyContainer) FindItemAwareById(id schema.IdRef) (IItemAware, bool) {
	return nil, false
}

func (p *PropertyContainer) FindItemAwareByName(name string) (IItemAware, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	item, ok := p.items[name]
	if !ok {
		return nil, false
	}

	return item, true
}

func (p *PropertyContainer) PutItemAwareByName(name string, itemAware IItemAware) {
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

func (p *PropertyContainer) CloneFor(other ILocatorCloner) {
	out, ok := other.(*PropertyContainer)
	if !ok {
		return
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	out.mu.Lock()
	defer out.mu.Unlock()

	cloneItemAwareMap(p.items, &out.items)
}

type FlowDataLocator struct {
	lmu       sync.RWMutex
	locators  map[string]IItemAwareLocator
	vmu       sync.RWMutex
	variables map[string]IItem
}

func NewFlowDataLocator() *FlowDataLocator {
	f := &FlowDataLocator{
		locators:  map[string]IItemAwareLocator{},
		variables: map[string]IItem{},
	}
	return f
}

func ElementToLocator(locator IFlowDataLocator, idGenerator id.IGenerator, element schema.Element) (err error) {
	dataObjectContainer := NewDataObjectContainer()
	if impl, ok := element.(interface {
		DataObjects() *[]schema.DataObject
	}); ok {
		for i := range *impl.DataObjects() {
			dataObject := &(*impl.DataObjects())[i]
			var name string
			if namePtr, present := dataObject.Name(); present {
				name = *namePtr
			} else {
				name = idGenerator.New().String()
			}
			container := NewContainer(dataObject)
			dataObjectBody := map[string]any{}
			if extension := dataObject.ExtensionElementsField; extension != nil {
				if extension.DataObjectBody != nil {
					err = json.Unmarshal([]byte(extension.DataObjectBody.Body), &dataObjectBody)
					if err != nil {
						err = fmt.Errorf("json Unmarshal DataObject %s: %v", *dataObject.IdField, err)
						return
					}
				}
			}
			container.Put(dataObjectBody)
			dataObjectContainer.dataObjectsByName[name] = container
			if idPtr, present := dataObject.Id(); present {
				dataObjectContainer.dataObjects[*idPtr] = container
			}
		}
	}

	if impl, ok := element.(interface {
		DataObjectReferences() *[]schema.DataObjectReference
	}); ok {
		for i := range *impl.DataObjectReferences() {
			dataObjectReference := &(*impl.DataObjectReferences())[i]
			var name string
			if namePtr, present := dataObjectReference.Name(); present {
				name = *namePtr
			} else {
				name = idGenerator.New().String()
			}
			var container IItemAware
			if dataObjPtr, present := dataObjectReference.DataObjectRef(); present {
				for dataObjectId := range dataObjectContainer.dataObjects {
					if dataObjectId == *dataObjPtr {
						container = dataObjectContainer.dataObjects[dataObjectId]
						break
					}
				}
				if container == nil {
					err = errors.NotFoundError{
						Expected: fmt.Sprintf("data object with ID %s", *dataObjPtr),
					}
					return
				}
			} else {
				err = errors.InvalidArgumentError{
					Expected: "data object reference to have dataObjectRef",
					Actual:   dataObjectReference,
				}
				return
			}
			dataObjectContainer.dataObjectReferencesByName[name] = container
			if idPtr, present := dataObjectReference.Id(); present {
				dataObjectContainer.dataObjectReferences[*idPtr] = container
			}
		}
	}

	if impl, ok := element.(interface {
		Properties() *[]schema.Property
	}); ok {
		for i := range *impl.Properties() {
			property := &(*impl.Properties())[i]
			var name string
			if namePtr, present := property.Name(); present {
				name = *namePtr
			} else {
				name = idGenerator.New().String()
			}
			container := NewContainer(property)
			dataObjectContainer.propertiesByName[name] = container
			if idPtr, present := property.Id(); present {
				dataObjectContainer.properties[*idPtr] = container
			}
		}
	}

	dol, found := locator.FindIItemAwareLocator(LocatorObject)
	if found {
		cloneFor, ok := dol.(ILocatorCloner)
		if ok {
			cloneFor.CloneFor(dataObjectContainer)
		}
	}
	locator.PutIItemAwareLocator(LocatorObject, dataObjectContainer)

	headerContainer, found := locator.FindIItemAwareLocator(LocatorHeader)
	if !found {
		headerContainer = NewHeaderContainer()
		locator.PutIItemAwareLocator(LocatorHeader, headerContainer)
	}
	propertyContainer, found := locator.FindIItemAwareLocator(LocatorProperty)
	if !found {
		propertyContainer = NewPropertyContainer()
		locator.PutIItemAwareLocator(LocatorProperty, propertyContainer)
	}
	if impl, ok := element.(schema.BaseElementInterface); ok {
		extensionElements, present := impl.ExtensionElements()
		if present {
			if headers := extensionElements.TaskHeaderField; headers != nil {
				for _, item := range headers.Header {
					container := NewContainer(nil)
					container.Put(item.ValueFor())
					headerContainer.PutItemAwareByName(item.Name, container)
				}
			}
			if properties := extensionElements.PropertiesField; properties != nil {
				for _, property := range properties.Property {
					container := NewContainer(nil)
					container.Put(property.ValueFor())
					propertyContainer.PutItemAwareByName(property.Name, container)
				}
			}
		}
	}

	return
}

func (f *FlowDataLocator) Merge(other IFlowDataLocator) {
	for key, value := range other.CloneVariables() {
		f.SetVariable(key, value)
	}

	for _, name := range []string{LocatorObject, LocatorHeader, LocatorProperty} {
		if out, ok := other.FindIItemAwareLocator(name); ok {
			in, ok := f.FindIItemAwareLocator(name)
			if ok {
				impl1, ok1 := in.(ILocatorCloner)
				impl2, ok2 := out.(ILocatorCloner)
				if ok1 && ok2 {
					impl2.CloneFor(impl1)
				}
			} else {
				f.PutIItemAwareLocator(name, out)
			}
		}
	}
}

func (f *FlowDataLocator) FindIItemAwareLocator(name string) (locator IItemAwareLocator, found bool) {
	f.lmu.RLock()
	defer f.lmu.RUnlock()
	locator, found = f.locators[name]
	return locator, found
}

func (f *FlowDataLocator) PutIItemAwareLocator(name string, locator IItemAwareLocator) {
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

func cloneItemAwareMap(in map[string]IItemAware, out *map[string]IItemAware) {
	for name, item := range in {
		outMap, ok := (*out)[name]
		if !ok {
			outMap = NewContainer(nil)
			(*out)[name] = outMap
		}
		impl1, ok1 := item.(ILocatorCloner)
		impl2, ok2 := outMap.(ILocatorCloner)
		if ok1 && ok2 {
			impl1.CloneFor(impl2)
		}
	}
}
