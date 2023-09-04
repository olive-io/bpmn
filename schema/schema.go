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

package schema

import (
	"encoding/xml"
	"reflect"
	"strconv"
	"strings"

	json "github.com/json-iterator/go"
)

// Base types

// QName XML qualified name (http://books.xmlschemata.org/relaxng/ch19-77287.html)
type QName string

func (t *QName) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	out := QName(*t)
	start.Name = xml.Name{
		Local: "bpmn:" + start.Name.Local,
	}

	return e.EncodeElement(out, start)
}

// Id Identifier (http://books.xmlschemata.org/relaxng/ch19-77151.html)
type Id = string

// IdRef Reference to identifiers (http://books.xmlschemata.org/relaxng/ch19-77159.html)
type IdRef = string

// AnyURI Corresponds normatively to the XLink href attribute (http://books.xmlschemata.org/relaxng/ch19-77009.html)
type AnyURI = string

// Element is implemented by every BPMN document element.
type Element interface {
	// FindBy recursively searches for an Element matching a predicate.
	//
	// Returns matching predicate in `result` and sets `found` to `true`
	// if any found.
	FindBy(ElementPredicate) (result Element, found bool)
}

type ElementPredicate func(Element) bool

// ExactId Returns a function that matches Element's identifier (if the `Element`
// implements BaseElementInterface against given string. If it matches, the
// function returns `true`.
//
// To be used in conjunction with FindBy (Element interface)
func ExactId(s string) ElementPredicate {
	return func(e Element) bool {
		if el, ok := e.(BaseElementInterface); ok {
			if id, present := el.Id(); present {
				return *id == s
			} else {
				return false
			}
		} else {
			return false
		}
	}
}

// ElementType Returns a function that matches Elements by types. If it matches, the
// function returns `true.
//
// To be used in conjunction with FindBy (Element interface)
func ElementType(t Element) ElementPredicate {
	return func(e Element) bool {
		return reflect.TypeOf(e) == reflect.TypeOf(t)
	}
}

// ElementInterface Returns a function that matches Elements by interface implementation. If it matches, the
// function returns `true.
//
// To be used in conjunction with FindBy (Element interface)
func ElementInterface(t interface{}) ElementPredicate {
	return func(e Element) bool {
		return reflect.TypeOf(e).Implements(reflect.TypeOf(t).Elem())
	}
}

// And Returns a function that matches Elements if both the receiver and the given
// predicates match.
//
// To be used in conjunction with FindBy (Element interface)
func (matcher ElementPredicate) And(andMatcher ElementPredicate) ElementPredicate {
	return func(e Element) bool {
		return matcher(e) && andMatcher(e)
	}
}

// Or Returns a function that matches Elements if either the receiver or the given
// predicates match.
//
// To be used in conjunction with FindBy (Element interface)
func (matcher ElementPredicate) Or(andMatcher ElementPredicate) ElementPredicate {
	return func(e Element) bool {
		return matcher(e) || andMatcher(e)
	}
}

// Special case for handling expressions being substituted for formal expression
// as per "BPMN 2.0 by Example":
//
// ```
//   <conditionExpression xsi:type="tFormalExpression">
//     ${getDataObject("TicketDataObject").status == "Resolved"}
//   </conditionExpression>
// ```
//
// Technically speaking, any type can be "patched" this way, but it seems impractical to
// do so broadly. Even within the confines of expressions, we don't want to support changing
// the type to an arbitrary one. Only specifying expressions as formal expressions will be
// possible for practicality.

// For this, a special type called AnExpression will be added and schema generator will
// replace Expression type with it.

// AnExpression Expression family container
//
// Expression field can be type-switched between Expression and FormalExpression
type AnExpression struct {
	Expression ExpressionInterface
}

func (e *AnExpression) MarshalXML(en *xml.Encoder, start xml.StartElement) (err error) {
	switch tt := e.Expression.(type) {
	case *FormalExpression:
		start.Attr = append(start.Attr, xml.Attr{
			Name: xml.Name{
				Local: "xsi:type",
			},
			Value: "tFormalExpression",
		})

		return en.EncodeElement(tt, start)
	case *Expression:
		start.Attr = append(start.Attr, xml.Attr{
			Name: xml.Name{
				Local: "xsi:type",
			},
			Value: "tExpression",
		})
		return en.EncodeElement(tt, start)
	}

	return
}

func (e *AnExpression) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	formal := false
	for i := range start.Attr {
		if start.Attr[i].Name.Space == "http://www.w3.org/2001/XMLSchema-instance" &&
			start.Attr[i].Name.Local == "type" &&
			// here we're checked for suffix instead of equality because
			// there doesn't seem to be a way to get a mapping between
			// bpmn schema and its namespace, so `bpmn:tFormalExpression`
			// equality check will fail if a different namespace name will
			// be used.
			(strings.HasSuffix(start.Attr[i].Value, ":tFormalExpression") ||
				start.Attr[i].Value == "tFormalExpression") {
			formal = true
			break
		}
	}
	if !formal {
		expr := DefaultExpression()
		err = d.DecodeElement(&expr, &start)
		if err != nil {
			return
		}
		e.Expression = &expr
	} else {
		expr := DefaultFormalExpression()
		err = d.DecodeElement(&expr, &start)
		if err != nil {
			return
		}
		e.Expression = &expr
	}
	return
}

func (e *AnExpression) FindBy(pred ElementPredicate) (result Element, found bool) {
	if e.Expression == nil {
		found = false
		return
	}
	return e.Expression.FindBy(pred)
}

// Equal will return true if both given elements are exactly the same ones
func Equal(e1, e2 Element) bool {
	// I'm not particularly fond of reflect.DeepEqual usage here,
	// perhaps we can come up with something better, perhaps a better
	// way to uniquely identify every element and compare that?
	return reflect.DeepEqual(e1, e2)
}

const Namespace = "bpmn:"

func PreMarshal(element Element, encoder *xml.Encoder, start *xml.StartElement) {
	if start.Name.Space == "http://www.omg.org/spec/BPMN/20100524/MODEL" {
		start.Name.Space = ""
		start.Name.Local = "bpmn:" + start.Name.Local
	}
	if start.Name.Space == "http://olive.io/spec/BPMN/MODEL" {
		start.Name.Space = ""
		start.Name.Local = "olive:" + start.Name.Local
	}

	if _, ok := element.(*Definitions); ok {
		start.Name.Local = Namespace + "definitions"
		start.Attr = append(start.Attr,
			xml.Attr{
				Name:  xml.Name{Local: "xmlns:bpmn"},
				Value: "http://www.omg.org/spec/BPMN/20100524/MODEL",
			},
			xml.Attr{
				Name:  xml.Name{Local: "xmlns:bpmndi"},
				Value: "http://www.omg.org/spec/BPMN/20100524/DI",
			},
			xml.Attr{
				Name:  xml.Name{Local: "xmlns:dc"},
				Value: "http://www.omg.org/spec/DD/20100524/DC",
			},
			xml.Attr{
				Name:  xml.Name{Local: "xmlns:di"},
				Value: "http://www.omg.org/spec/DD/20100524/DI",
			},
			xml.Attr{
				Name:  xml.Name{Local: "xmlns:olive"},
				Value: "http://olive.io/spec/BPMN/MODEL",
			},
		)
	}
}

type ItemType string

const (
	ItemTypeObject  ItemType = "object"
	ItemTypeInteger ItemType = "integer"
	ItemTypeString  ItemType = "string"
	ItemTypeBoolean ItemType = "boolean"
	ItemTypeFloat   ItemType = "float"
)

type TaskDefinition struct {
	Type string `xml:"type,attr"`
}

func (t *TaskDefinition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	out := TaskDefinition(*t)
	start.Name = xml.Name{
		Local: "olive:taskDefinition",
	}

	return e.EncodeElement(out, start)
}

type Item struct {
	Key   string   `xml:"key,attr"`
	Value string   `xml:"value,attr"`
	Type  ItemType `xml:"type,attr"`
}

func (i *Item) ValueFor() any {
	switch i.Type {
	case ItemTypeString:
		return i.Value
	case ItemTypeInteger:
		integer, _ := strconv.ParseInt(i.Value, 10, 64)
		return int(integer)
	case ItemTypeBoolean:
		return i.Value == "true"
	case ItemTypeFloat:
		f, _ := strconv.ParseFloat(i.Value, 64)
		return f
	case ItemTypeObject:
		obj := map[string]interface{}{}
		_ = json.Unmarshal([]byte(i.Value), &obj)
		return obj
	default:
		return nil
	}
}

func (i *Item) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	out := Item(*i)
	start.Name = xml.Name{
		Local: "olive:item",
	}

	return e.EncodeElement(out, start)
}

type TaskHeader struct {
	ItemFields []*Item `xml:"http://olive.io/spec/BPMN/MODEL item"`
}

func (h *TaskHeader) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	out := TaskHeader(*h)
	start.Name = xml.Name{
		Local: "olive:taskHeader",
	}

	return e.EncodeElement(out, start)
}

type Properties struct {
	ItemFields []*Item `xml:"http://olive.io/spec/BPMN/MODEL item"`
}

func (p *Properties) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	out := Properties(*p)
	start.Name = xml.Name{
		Local: "olive:properties",
	}

	return e.EncodeElement(out, start)
}

// Generate schema files:

//go:generate saxon-he ../BPMN20.xsd ../schema-codegen.xsl
