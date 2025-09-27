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

package expr

import (
	"context"
	"testing"

	"github.com/expr-lang/expr/ast"
	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/expression"
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

func TestExprSum(t *testing.T) {
	var engine expression.IEngine = New(context.Background())
	compiled, err := engine.CompileExpression("a + b")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, map[string]interface{}{
		"a": 1,
		"b": 2,
	})
	assert.Nil(t, err)
	_, ok := result.(int)
	assert.True(t, ok)
}

type Visitor struct {
	Identifiers []string
}

func (v *Visitor) Visit(node *ast.Node) {
	if n, ok := (*node).(*ast.IdentifierNode); ok {
		v.Identifiers = append(v.Identifiers, n.Value)
	}
}

type dataObjects map[string]data.IItemAware

func (d dataObjects) PutItemAwareById(id schema.IdRef, itemAware data.IItemAware) {
}

func (d dataObjects) PutItemAwareByName(name string, itemAware data.IItemAware) {
}

func (d dataObjects) FindItemAwareById(id schema.IdRef) (itemAware data.IItemAware, found bool) {
	itemAware, found = d[id]
	return
}

func (d dataObjects) FindItemAwareByName(name string) (itemAware data.IItemAware, found bool) {
	itemAware, found = d[name]
	return
}

func (d dataObjects) Clone() map[string]data.IItem {
	return map[string]data.IItem{}
}

func TestExpr_getDataObject(t *testing.T) {
	var engine = New(context.Background())
	container := data.NewContainer(nil)
	container.Put(schema.NewValue(1))
	container1 := data.NewContainer(nil)
	container1.Put(schema.NewValue(map[string]string{"msg": "hello"}))
	var objs dataObjects = map[string]data.IItemAware{
		"dataObject":  container,
		"dataObject1": container1,
	}
	engine.SetItemAwareLocator(data.LocatorObject, objs)

	h1 := data.NewContainer(nil)
	h1.Put(schema.NewValue(2))
	var headers dataObjects = map[string]data.IItemAware{
		"a": h1,
	}
	engine.SetItemAwareLocator(data.LocatorHeader, headers)

	p1 := data.NewContainer(nil)
	p1.Put(schema.NewValue("bar"))
	var properties dataObjects = map[string]data.IItemAware{
		"foo": p1,
	}
	engine.SetItemAwareLocator(data.LocatorProperty, properties)

	var properties1 = map[string]any{
		"foo": "bar",
	}

	compiled, err := engine.CompileExpression("getDataObject('dataObject1').msg == 'hello'")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, properties1)
	assert.Nil(t, err)
	assert.True(t, result.(bool))
	compiled, err = engine.CompileExpression("getDataObject('dataObject') == 1")
	assert.Nil(t, err)
	result, err = engine.EvaluateExpression(compiled, properties1)
	assert.Nil(t, err)
	assert.True(t, result.(bool))

	compiled, err = engine.CompileExpression("getHeader('a') + 1")
	assert.Nil(t, err)
	result, err = engine.EvaluateExpression(compiled, properties1)
	assert.Nil(t, err)
	assert.Equal(t, result.(int), int(3))

	compiled, err = engine.CompileExpression("getHeader('a') < 0")
	assert.Nil(t, err)
	result, err = engine.EvaluateExpression(compiled, properties1)
	assert.Nil(t, err)
	assert.Equal(t, result.(bool), false)

	compiled, err = engine.CompileExpression("getProp('foo') == 'bar' or getProp('foo') != 'bar'")
	assert.Nil(t, err)
	result, err = engine.EvaluateExpression(compiled, properties1)
	assert.Nil(t, err)
	assert.True(t, result.(bool))

	compiled, err = engine.CompileExpression("foo == 'bar'")
	assert.Nil(t, err)
	result, err = engine.EvaluateExpression(compiled, properties1)
	assert.Nil(t, err)
	assert.True(t, result.(bool))
}
