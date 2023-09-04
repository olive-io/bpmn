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

package data

import (
	"context"

	"github.com/olive-io/bpmn/schema"
)

// IItem is an abstract interface for a piece of data
type IItem interface{}

// IIteratorStopper stops Collection iterator and releases resources
// associated with it
type IIteratorStopper interface {
	// Stop does the actual stopping
	Stop()
}

type channelIteratorStopper struct {
	ch chan struct{}
}

func makeChannelIteratorStopper() channelIteratorStopper {
	return channelIteratorStopper{ch: make(chan struct{})}
}

func (c channelIteratorStopper) close() {
	close(c.ch)
}

func (c channelIteratorStopper) Stop() {
	c.ch <- struct{}{}
}

type ICollection interface {
	IItem
	// ItemIterator returns a channel that iterates over collection's
	// items and an IteratorStopper that must be used if iterator was
	// not exhausted, otherwise there'll be a memory leak in a form
	// of a goroutine that does nothing.
	//
	// The iterator will also clean itself up and terminate upon
	// context termination.
	ItemIterator(ctx context.Context) (chan IItem, IIteratorStopper)
}

type SliceIterator []IItem

func (s *SliceIterator) ItemIterator(ctx context.Context) (items chan IItem, stop IIteratorStopper) {
	items = make(chan IItem)
	stopper := makeChannelIteratorStopper()
	stop = stopper
	go func() {
	loop:
		for i := range *s {
			select {
			case <-ctx.Done():
				break loop
			case <-stopper.ch:
				break loop
			case items <- (*s)[i]:
			}
		}
		close(items)
		stopper.close()
	}()
	return
}

// ItemOrCollection will return nil if no items given,
// the same item if only one item is given and SliceIterator
// if more than one item is given. SliceIterator implements
// Collection and, therefore, also implements Item.
func ItemOrCollection(items ...IItem) (item IItem) {
	switch len(items) {
	case 0:
	case 1:
		item = items[0]
	default:
		item = SliceIterator(items)
	}
	return
}

// IItemAware provides basic interface of accessing data items
type IItemAware interface {
	// Get returns a channel that will eventually return the data item
	Get() IItem
	// Put puts to update the item
	Put(item IItem)
}

// IItemAwareLocator interface describes a way to find and put IItemAware
type IItemAwareLocator interface {
	// FindItemAwareById finds ItemAware by its schema.Id
	FindItemAwareById(id schema.IdRef) (itemAware IItemAware, found bool)
	// FindItemAwareByName finds ItemAware by its name (where applicable)
	FindItemAwareByName(name string) (itemAware IItemAware, found bool)
	// Clone clones all IItem to the specified target
	Clone() map[string]any
}

type DefaultItemAwareLocator struct{}

func (d DefaultItemAwareLocator) FindItemAwareById(id schema.IdRef) (itemAware IItemAware, found bool) {
	return
}

func (d DefaultItemAwareLocator) FindItemAwareByName(name string) (itemAware IItemAware, found bool) {
	return
}

func (d DefaultItemAwareLocator) Clone() map[string]any { return map[string]any{} }
