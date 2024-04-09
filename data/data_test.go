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

package data

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceIterator_ItemIterator(t *testing.T) {
	items := SliceIterator([]IItem{1, "hello"})
	ctx := context.Background()
	iterator, _ := items.ItemIterator(ctx)
	newItems := SliceIterator(make([]IItem, 0, len(items)))
	for e := range iterator {
		newItems = append(newItems, e)
	}
	assert.Equal(t, items, newItems)
}

func TestSliceIterator_ItemIterator_Stopping(t *testing.T) {
	items := SliceIterator([]IItem{1, "hello"})
	ctx := context.Background()
	iterator, stopper := items.ItemIterator(ctx)
	stopper.Stop()
	value, ok := <-iterator
	assert.Nil(t, value)
	assert.False(t, ok)
}

func TestSliceIterator_ItemIterator_ContextDone(t *testing.T) {
	items := SliceIterator([]IItem{1, "hello"})
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
