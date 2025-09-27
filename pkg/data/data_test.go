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

package data

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
)

func TestSliceIterator_ItemIterator(t *testing.T) {
	items := NewSlice([]IItem{schema.NewValue(1), schema.NewValue("hello")})
	ctx := context.Background()
	iterator, _ := items.ItemIterator(ctx)
	newItems := make([]IItem, 0)
	for e := range iterator {
		newItems = append(newItems, e)
	}
	assert.Equal(t, items.items, newItems)
}

func TestSliceIterator_ItemIterator_Stopping(t *testing.T) {
	items := NewSlice([]IItem{schema.NewValue(1), schema.NewValue("hello")})
	ctx := context.Background()
	iterator, stopper := items.ItemIterator(ctx)
	stopper.Stop()
	value, ok := <-iterator
	assert.Nil(t, value)
	assert.False(t, ok)
}

func TestSliceIterator_ItemIterator_ContextDone(t *testing.T) {
	items := NewSlice([]IItem{schema.NewValue(1), schema.NewValue("hello")})
	ctx, cancel := context.WithCancel(context.Background())
	iterator, stop := items.ItemIterator(ctx)
	cancel()

	// Wait until stopper is closed. This is definitely not
	// black box testing, but this essentially ensures
	// the cancellation has been handled
	stopper := stop.(channelIteratorStopper)
	<-stopper.ch

	value, ok := <-iterator
	assert.Nil(t, value)
	assert.False(t, ok)
}
