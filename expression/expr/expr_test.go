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
	compiled, err := engine.CompileExpression("$('dataObject1').msg == 'hello'")
	assert.Nil(t, err)
	compiled, err = engine.CompileExpression("$('dataObject') == 1")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, map[string]interface{}{})
	assert.Nil(t, err)
	assert.True(t, result.(bool))
}
