/*
   Copyright 2023 The bpmn Authors

   This library is free software; you can redistribute it and/or
   modify it under the terms of the GNU Lesser General Public
   License as published by the Free Software Foundation; either
   version 2.1 of the License, or (at your option) any later version.

   This library is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
   Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public
   License along with this library;
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

func (d dataObjects) Clone() map[string]any {
	return map[string]any{}
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
	engine.SetItemAwareLocator(data.LocatorObject, objs)

	h1 := data.NewContainer(nil)
	h1.Put(2)
	var headers dataObjects = map[string]data.IItemAware{
		"a": h1,
	}
	engine.SetItemAwareLocator(data.LocatorHeader, headers)

	compiled, err := engine.CompileExpression("d$('dataObject1').msg == 'hello'")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, map[string]interface{}{})
	assert.Nil(t, err)
	assert.True(t, result.(bool))
	compiled, err = engine.CompileExpression("d$('dataObject') == 1")
	assert.Nil(t, err)
	result, err = engine.EvaluateExpression(compiled, map[string]interface{}{})
	assert.Nil(t, err)
	assert.True(t, result.(bool))

	compiled, err = engine.CompileExpression("h$('a') + 1")
	assert.Nil(t, err)
	result, err = engine.EvaluateExpression(compiled, map[string]interface{}{})
	assert.Nil(t, err)
	assert.Equal(t, result.(int), int(3))
}
