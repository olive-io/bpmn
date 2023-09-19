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
