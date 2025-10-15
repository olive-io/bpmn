/*
Copyright 2024 The bpmn Authors

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

package data

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
)

type DataItem struct {
	value string
}

func (di *DataItem) Type() schema.ItemType {
	return schema.ItemTypeObject
}

func (di *DataItem) Value() any {
	return map[string]interface{}{"name": di.value}
}

func TestFlowDataLocator_Merge(t *testing.T) {
	l1 := NewFlowDataLocator()
	l2 := NewFlowDataLocator()
	l1.SetVariable("a", "b")
	l1.SetVariable("c", &DataItem{value: "dd"})
	l2.Merge(l1)

	v1 := l1.CloneVariables()["a"]
	v2 := l2.CloneVariables()["a"]
	assert.Equal(t, v1.Value(), v2.Value())
	vv, _ := l1.GetVariable("c")
	assert.Equal(t, vv, map[string]any{"name": "dd"})

	locator := NewPropertyContainer()
	c1 := NewContainer(nil)
	c1.Put(schema.NewValue("aa"))
	locator.PutItemAwareByName("a", c1)
	l1.PutIItemAwareLocator(LocatorProperty, locator)

	l2.Merge(l1)

	l2Locator, _ := l2.FindIItemAwareLocator(LocatorProperty)
	cloned := l2Locator.Clone()
	value := cloned["a"].Value()
	assert.Equal(t, value, "aa")
}

func TestFlowDataLocator_CloneItems(t *testing.T) {
	l := NewFlowDataLocator()
	aware := NewContainer(nil)
	aware.Put(schema.NewValue("hello"))
	container := NewDataObjectContainer()
	container.PutItemAwareById("id", aware)
	container.PutItemAwareByName("in", aware)
	l.PutIItemAwareLocator(LocatorObject, container)

	containerOut, ok := l.FindIItemAwareLocator(LocatorObject)
	if !assert.True(t, ok) {
		return
	}

	if !assert.Equal(t, containerOut, container) {
		return
	}

	awareOut, ok := container.FindItemAwareById("id")
	if !assert.True(t, ok) {
		return
	}

	if !assert.Equal(t, awareOut, aware) {
		return
	}

	awareOut.Put(schema.NewValue("hello"))

	items := l.CloneItems(LocatorObject)
	if !assert.Equal(t, items["id"], schema.NewValue("hello")) {
		return
	}
}
