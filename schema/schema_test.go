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

package schema

import (
	"fmt"
	"testing"

	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestItem_ValueTo(t *testing.T) {
	type fields struct {
		Name  string
		Value string
		Type  ItemType
	}
	type args struct {
		dst any
	}
	intDst := int64(1)
	uintDst := uint32(11)
	floatDst := 2.2
	arrDist := make([]any, 0)
	mapDist := map[string]any{}

	type SchemaData struct {
		Name string   `json:"name"`
		Type ItemType `json:"type"`
	}

	structDist := SchemaData{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"valueTo_int", fields{"", "1", ItemTypeInteger}, args{&intDst}, assert.Error},
		{"valueTo_uint", fields{"", "11", ItemTypeInteger}, args{&uintDst}, assert.Error},
		{"valueTo_float", fields{"", "2.2", ItemTypeFloat}, args{&floatDst}, assert.Error},
		{"valueTo_slice", fields{"", `[1, 1.1, "hello"]`, ItemTypeArray}, args{&arrDist}, assert.Error},
		{"valueTo_map", fields{"", `{"a": 1, "b": "bb", "name": "scott"}`, ItemTypeObject}, args{&mapDist}, assert.Error},
		{"valueTo_struct", fields{"", `{"name": "Task", "type": "string"}`, ItemTypeObject}, args{&structDist}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
				Type:  tt.fields.Type,
			}
			tt.wantErr(t, i.ValueTo(tt.args.dst), fmt.Sprintf("ValueTo(%v)", tt.args.dst))
			data, _ := json.Marshal(tt.args.dst)
			t.Logf("v = %v\n", string(data))
		})
	}
}
