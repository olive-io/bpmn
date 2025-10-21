/*
Copyright 2025 The bpmn Authors

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
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"

	json "github.com/bytedance/sonic"
)

type ItemType string

func (it ItemType) String() string { return string(it) }

const (
	ItemTypeObject  ItemType = "object"
	ItemTypeArray   ItemType = "array"
	ItemTypeInteger ItemType = "integer"
	ItemTypeString  ItemType = "string"
	ItemTypeBoolean ItemType = "boolean"
	ItemTypeFloat   ItemType = "float"
)

type Value struct {
	ItemType  ItemType `json:"type" xml:"type,attr"`
	ItemValue string   `json:"value" xml:"value,attr"`
}

func (iv *Value) Type() ItemType { return iv.ItemType }

func (iv *Value) Value() any { return iv.ValueFor() }

func NewValue(value any) *Value {
	v := &Value{}
	v.ValueFrom(value)
	return v
}

func (iv *Value) ValueFor() any {
	switch iv.Type() {
	case ItemTypeString:
		return iv.ItemValue
	case ItemTypeInteger:
		n, _ := strconv.ParseInt(iv.ItemValue, 10, 64)
		return n
	case ItemTypeBoolean:
		return iv.ItemValue == "true"
	case ItemTypeFloat:
		f, _ := strconv.ParseFloat(iv.ItemValue, 64)
		return f
	case ItemTypeArray:
		var arr []any
		_ = json.Unmarshal([]byte(iv.ItemValue), &arr)
		return arr
	case ItemTypeObject:
		obj := map[string]any{}
		_ = json.Unmarshal([]byte(iv.ItemValue), &obj)
		return obj
	default:
		return iv.ItemValue
	}
}

func (iv *Value) ValueFrom(value any) {
	vv, ok := value.(*Value)
	if ok {
		iv.ItemType = vv.ItemType
		iv.ItemValue = vv.ItemValue
		return
	}

	switch iv.ItemType {
	case ItemTypeString:
		v, ok := value.(string)
		if ok {
			iv.ItemValue = v
		}
	case ItemTypeInteger:
		switch v := value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			iv.ItemValue = fmt.Sprintf("%v", v)
		case string:
			_, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return
			}
			iv.ItemValue = v
		}
	case ItemTypeBoolean:
		v, ok := value.(bool)
		if ok {
			iv.ItemValue = strconv.FormatBool(v)
		} else {
			if b, ok := value.(string); ok {
				if b == "true" || b == "false" {
					iv.ItemValue = b
				}
			}
		}
	case ItemTypeFloat:
		switch vv := value.(type) {
		case float64, float32:
			iv.ItemValue = fmt.Sprintf("%f", vv)
		case string:
			_, err := strconv.ParseFloat(vv, 64)
			if err != nil {
				return
			}
			iv.ItemValue = vv
		}
	case ItemTypeArray:
		rt := reflect.TypeOf(value)
		if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
			data, err := json.Marshal(value)
			if err != nil {
				return
			}
			iv.ItemValue = string(data)
		}
		if rt.Kind() == reflect.String {
			vv := value.(string)
			var arr []any
			if err := json.Unmarshal([]byte(vv), &arr); err != nil {
				return
			}
			iv.ItemValue = vv
		}
	case ItemTypeObject:
		rt := reflect.TypeOf(value)
		if rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}
		if rt.Kind() == reflect.Struct || rt.Kind() == reflect.Map {
			data, err := json.Marshal(value)
			if err != nil {
				return
			}
			iv.ItemValue = string(data)
		}
	default:
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}

		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			iv.ItemType = ItemTypeArray
			data, err := json.Marshal(value)
			if err != nil {
				return
			}
			iv.ItemValue = string(data)
		case reflect.Map, reflect.Struct:
			iv.ItemType = ItemTypeObject
			data, err := json.Marshal(value)
			if err != nil {
				return
			}
			iv.ItemValue = string(data)
		case reflect.String:
			iv.ItemType = ItemTypeString
			iv.ItemValue = rv.String()
		case reflect.Bool:
			iv.ItemType = ItemTypeBoolean
			iv.ItemValue = strconv.FormatBool(rv.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			iv.ItemType = ItemTypeInteger
			iv.ItemValue = fmt.Sprintf("%v", rv.Int())
		case reflect.Float32, reflect.Float64:
			iv.ItemType = ItemTypeFloat
			iv.ItemValue = fmt.Sprintf("%v", rv.Float())
		default:
		}

	}
}

func (iv *Value) ValueTo(dst any) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.String && iv.ItemType == ItemTypeString {
		rv.SetString(iv.ItemValue)
	} else if rv.Kind() == reflect.Map && iv.ItemType == ItemTypeObject {
		if err := json.Unmarshal([]byte(iv.ItemValue), &dst); err != nil {
			return err
		}
	} else if rv.Kind() == reflect.Struct && iv.ItemType == ItemTypeObject {
		if err := json.Unmarshal([]byte(iv.ItemValue), &dst); err != nil {
			return err
		}
	} else if rv.Kind() == reflect.Slice && iv.ItemType == ItemTypeArray {
		if err := json.Unmarshal([]byte(iv.ItemValue), &dst); err != nil {
			return err
		}
	} else if rv.CanInt() && iv.ItemType == ItemTypeInteger {
		n, err := strconv.ParseInt(iv.ItemValue, 10, 64)
		if err != nil {
			return err
		}
		rv.SetInt(n)
	} else if rv.CanUint() && iv.ItemType == ItemTypeInteger {
		n, err := strconv.ParseUint(iv.ItemValue, 10, 64)
		if err != nil {
			return err
		}
		rv.SetUint(n)
	} else if rv.CanFloat() && iv.ItemType == ItemTypeFloat {
		f, err := strconv.ParseFloat(iv.ItemValue, 10)
		if err != nil {
			return err
		}
		rv.SetFloat(f)
	}

	return fmt.Errorf("not matched")
}

func (iv *Value) DeepCopy() *Value {
	v := &Value{
		ItemType:  iv.ItemType,
		ItemValue: iv.ItemValue,
	}
	return v
}

type Item struct {
	Name  string   `xml:"name,attr"`
	Value string   `xml:"value,attr"`
	Type  ItemType `xml:"type,attr"`
}

func (i *Item) ToValue() *Value {
	return &Value{ItemType: i.Type, ItemValue: i.Value}
}

func (i *Item) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	out := Item(*i)
	start.Name = xml.Name{
		Local: OliveNS + start.Name.Local,
	}

	if out.Type == "" {
		out.Type = ItemTypeString
	}

	if err := e.EncodeElement(out, start); err != nil {
		return err
	}

	return nil
}

func (i *Item) UnmarshalXML(de *xml.Decoder, start xml.StartElement) error {
	type ItemUnmarshaler Item
	out := ItemUnmarshaler{}
	if err := de.DecodeElement(&out, &start); err != nil {
		return err
	}

	*i = Item(out)
	if i.Type == "" {
		i.Type = ItemTypeString
	}
	return nil
}
