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
	"sync"

	"github.com/olive-io/bpmn/schema"
)

type Container struct {
	schema.ItemAwareInterface
	mu   sync.RWMutex
	item IItem
}

func NewContainer(itemAware schema.ItemAwareInterface) *Container {
	container := &Container{
		ItemAwareInterface: itemAware,
	}
	return container
}

func (c *Container) Unavailable() bool {
	return false
}

func (c *Container) Get() IItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.item
}

func (c *Container) Put(item IItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.item = item
}

func (c *Container) CloneFor(other ILocatorCloner) {
	out, ok := other.(*Container)
	if !ok {
		return
	}
	if out.ItemAwareInterface == nil {
		out.ItemAwareInterface = c.ItemAwareInterface
	}
	value := c.Get()
	if value != nil {
		out.Put(value)
	}
}
