package expr

import (
	"context"
	"testing"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/expression"
	"github.com/olive-io/bpmn/schema"
	"github.com/stretchr/testify/assert"
)

func TestExpr(t *testing.T) {
	var engine expression.Engine = New(context.Background())
	compiled, err := engine.CompileExpression("a > 1")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, map[string]interface{}{
		"a": 2,
	})
	assert.Nil(t, err)
	assert.True(t, result.(bool))
}

type dataObjects map[string]data.ItemAware

func (d dataObjects) FindItemAwareById(id schema.IdRef) (itemAware data.ItemAware, found bool) {
	itemAware, found = d[id]
	return
}

func (d dataObjects) FindItemAwareByName(name string) (itemAware data.ItemAware, found bool) {
	itemAware, found = d[name]
	return
}

func TestExpr_getDataObject(t *testing.T) {
	var engine = New(context.Background())
	container := data.NewContainer(context.Background(), nil)
	container.Put(context.Background(), 1)
	var objs dataObjects = map[string]data.ItemAware{
		"dataObject": container,
	}
	engine.SetItemAwareLocator(objs)
	compiled, err := engine.CompileExpression("getDataObject('dataObject') > 0")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, map[string]interface{}{})
	assert.Nil(t, err)
	assert.True(t, result.(bool))
}
