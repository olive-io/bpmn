package xpath

import (
	"context"
	"testing"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/expression"
	"github.com/olive-io/bpmn/schema"
	"github.com/stretchr/testify/assert"
)

func TestXPath(t *testing.T) {
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

func (d dataObjects) FindItemAwareById(id schema.IdRef) (itemAware data.IItemAware, found bool) {
	itemAware, found = d[id]
	return
}

func (d dataObjects) FindItemAwareByName(name string) (itemAware data.IItemAware, found bool) {
	itemAware, found = d[name]
	return
}

func (d dataObjects) Clone() map[string]any {
	panic("implement me")
}

func TestXPath_getDataObject(t *testing.T) {
	// This funtionality doesn't quite work yet
	t.SkipNow()
	var engine = New(context.Background())
	container := data.NewContainer(nil)
	container.Put(data.XMLSource(`<tag attr="val"/>`))
	var objs dataObjects = map[string]data.IItemAware{
		"dataObject": container,
	}
	engine.SetItemAwareLocator("$", objs)
	compiled, err := engine.CompileExpression("(getDataObject('dataObject')/tag/@attr/string())[1]")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, map[string]interface{}{})
	assert.Nil(t, err)
	assert.Equal(t, "val", result.(string))
}
