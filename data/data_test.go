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
