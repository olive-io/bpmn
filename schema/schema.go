package schema

import (
	"encoding/xml"
	"io"
	"reflect"
	"strings"
)

// Base types

// QName XML qualified name (http://books.xmlschemata.org/relaxng/ch19-77287.html)
type QName = string

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
	if !strings.HasPrefix(start.Name.Local, Namespace) {
		start.Name.Local = Namespace + start.Name.Local
	}

	if _, ok := element.(FlowNodeInterface); ok {
		for _, attr := range start.Attr {
			_ = attr
		}
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

func (t *TaskDefinition) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "type" {
			t.Type = attr.Value
		}
	}

	return nil
}

type Item struct {
	Key   string   `xml:"key,attr"`
	Value string   `xml:"value,attr"`
	Type  ItemType `xml:"type,attr"`
}

func (i *Item) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	out := Item(*i)
	start.Name = xml.Name{
		Local: "olive:item",
	}

	return e.EncodeElement(out, start)
}

type TaskHeader struct {
	ItemFields []*Item
}

func (h *TaskHeader) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	out := TaskHeader(*h)
	start.Name = xml.Name{
		Local: "olive:taskHeader",
	}

	return e.EncodeElement(out, start)
}

func (h *TaskHeader) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	h.ItemFields = make([]*Item, 0)
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch tt := tok.(type) {
		case xml.StartElement:
			fname := strings.Join([]string{tt.Name.Space, tt.Name.Local}, ":")
			if fname == "olive:item" {
				item := &Item{}
				_ = d.DecodeElement(item, &tt)
				h.ItemFields = append(h.ItemFields, item)
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}
	}

	return nil
}

type Properties struct {
	ItemFields []*Item
}

func (p *Properties) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	out := Properties(*p)
	start.Name = xml.Name{
		Local: "olive:properties",
	}

	return e.EncodeElement(out, start)
}

func (p *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	p.ItemFields = make([]*Item, 0)
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch tt := tok.(type) {
		case xml.StartElement:
			fname := strings.Join([]string{tt.Name.Space, tt.Name.Local}, ":")
			if fname == "olive:item" {
				item := &Item{}
				_ = d.DecodeElement(item, &tt)
				p.ItemFields = append(p.ItemFields, item)
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}
	}

	return nil
}

func (t *ExtensionElements) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch tt := tok.(type) {
		case xml.StartElement:
			fname := strings.Join([]string{tt.Name.Space, tt.Name.Local}, ":")
			if fname == "olive:taskDefinition" {
				t.TaskDefinitionField = &TaskDefinition{}
				_ = t.TaskDefinitionField.UnmarshalXML(d, tt)
			}
			if fname == "olive:taskHeaders" {
				t.TaskHeaderField = &TaskHeader{}
				_ = t.TaskHeaderField.UnmarshalXML(d, tt)
			}
			if fname == "olive:properties" {
				t.PropertiesField = &Properties{}
				_ = t.PropertiesField.UnmarshalXML(d, tt)
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}

	}

	return nil
}

// Generate schema files:

//go:generate saxon-he ../BPMN20.xsd ../schema-codegen.xsl
