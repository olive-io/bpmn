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

package xpath

import (
	"context"
	"testing"

	"github.com/olive-io/bpmn/schema"
	"github.com/stretchr/testify/assert"

	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/expression"
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

func TestXPath_getDataObject(t *testing.T) {
	// This funtionality doesn't quite work yet
	t.SkipNow()
	var engine = New(context.Background())
	container := data.NewContainer(nil)
	container.Put(data.XMLSource(`<tag attr="val"/>`))
	var objs dataObjects = map[string]data.IItemAware{
		"dataObject": container,
	}
	engine.SetItemAwareLocator(data.LocatorObject, objs)
	compiled, err := engine.CompileExpression("($('dataObject')/tag/@attr/string())[1]")
	assert.Nil(t, err)
	result, err := engine.EvaluateExpression(compiled, map[string]interface{}{})
	assert.Nil(t, err)
	assert.Equal(t, "val", result.(string))
}
