/*
Copyright 2024 The bpmn Authors

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlowDataLocator_Merge(t *testing.T) {
	l1 := NewFlowDataLocator()
	l2 := NewFlowDataLocator()
	l1.SetVariable("a", "b")
	l2.Merge(l1)

	assert.Equal(t, l1.CloneVariables(), l2.CloneVariables())

	locator := NewPropertyContainer()
	c1 := NewContainer(nil)
	c1.Put("aa")
	locator.PutItemAwareByName("a", c1)
	l1.PutIItemAwareLocator(LocatorProperty, locator)

	l2.Merge(l1)

	l2Locator, _ := l2.FindIItemAwareLocator(LocatorProperty)
	assert.Equal(t, l2Locator.Clone()["a"], "aa")
}

func TestFlowDataLocator_CloneItems(t *testing.T) {
	l := NewFlowDataLocator()
	aware := NewContainer(nil)
	aware.Put("hello")
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

	awareOut.Put("hello")

	items := l.CloneItems(LocatorObject)
	if !assert.Equal(t, items["id"], "hello") {
		return
	}
}
