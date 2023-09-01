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
	var engine expression.IEngine = New(context.Background())
	compiled, err := engine.CompileExpression("a > 1")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, map[string]interface{}{
		"a": 2,
	})
	assert.Nil(t, err)
	assert.True(t, result.(bool))
}

type dataObjects map[string]data.IItemAware

func (d dataObjects) PutItemAwareById(id schema.IdRef) (itemAware data.IItemAware, found bool) {
	panic("implement me")
}

func (d dataObjects) PutItemAwareByName(name string) (itemAware data.IItemAware, found bool) {
	panic("implement me")
}

func (d dataObjects) FindItemAwareById(id schema.IdRef) (itemAware data.IItemAware, found bool) {
	itemAware, found = d[id]
	return
}

func (d dataObjects) FindItemAwareByName(name string) (itemAware data.IItemAware, found bool) {
	itemAware, found = d[name]
	return
}

func TestExpr_getDataObject(t *testing.T) {
	var engine = New(context.Background())
	container := data.NewContainer(nil)
	container.Put(1)
	container1 := data.NewContainer(nil)
	container1.Put(map[string]string{"msg": "hello"})
	var objs dataObjects = map[string]data.IItemAware{
		"dataObject":  container,
		"dataObject1": container1,
	}
	engine.SetItemAwareLocator("$", objs)
	compiled, err := engine.CompileExpression("$('dataObject1').msg == 'hello'")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, map[string]interface{}{})
	assert.Nil(t, err)
	assert.True(t, result.(bool))
}
