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
