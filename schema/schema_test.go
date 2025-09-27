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

package schema

import (
	"fmt"
	"testing"

	json "github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
)

func TestValue_ValueTo(t *testing.T) {
	type fields struct {
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
		{"valueTo_int", fields{"1", ItemTypeInteger}, args{&intDst}, assert.Error},
		{"valueTo_uint", fields{"11", ItemTypeInteger}, args{&uintDst}, assert.Error},
		{"valueTo_float", fields{"2.2", ItemTypeFloat}, args{&floatDst}, assert.Error},
		{"valueTo_slice", fields{`[1, 1.1, "hello"]`, ItemTypeArray}, args{&arrDist}, assert.Error},
		{"valueTo_map", fields{`{"a": 1, "b": "bb", "name": "scott"}`, ItemTypeObject}, args{&mapDist}, assert.Error},
		{"valueTo_struct", fields{`{"name": "Task", "type": "string"}`, ItemTypeObject}, args{&structDist}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Value{
				ItemValue: tt.fields.Value,
				ItemType:  tt.fields.Type,
			}
			tt.wantErr(t, i.ValueTo(tt.args.dst), fmt.Sprintf("ValueTo(%v)", tt.args.dst))
			data, _ := json.Marshal(tt.args.dst)
			t.Logf("v = %v\n", string(data))
		})
	}
}

func TestValue_ValueFrom(t *testing.T) {
	type fields struct {
		Value string
		Type  ItemType
	}
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ComparisonAssertionFunc
	}{
		{"valueFrom_int", fields{"", ItemTypeInteger}, args{value: int64(1)}, assert.Equal},
		{"valueFrom_int_with_error", fields{"", ItemTypeInteger}, args{value: "a"}, assert.Equal},
		{"valueFrom_string", fields{"", ItemTypeString}, args{value: "a"}, assert.Equal},
		{"valueFrom_string_1", fields{"", ""}, args{value: "a"}, assert.Equal},
		{"valueFrom_slice", fields{"", ItemTypeArray}, args{value: []int{1, 2, 3}}, assert.Equal},
		{"valueFrom_object", fields{"", ItemTypeObject}, args{value: map[string]int{"a": 1, "b": 2}}, assert.Equal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Value{
				ItemValue: tt.fields.Value,
				ItemType:  tt.fields.Type,
			}
			i.ValueFrom(tt.args.value)
			t.Logf("v = %v\n", i.ValueFor())
		})
	}
}
