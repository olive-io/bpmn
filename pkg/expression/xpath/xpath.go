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

package xpath

import (
	"bytes"
	"context"
	"reflect"

	"github.com/ChrisTrenkamp/xsel/exec"
	"github.com/ChrisTrenkamp/xsel/grammar"
	"github.com/ChrisTrenkamp/xsel/parser"
	"github.com/ChrisTrenkamp/xsel/store"
	"github.com/Chronokeeper/anyxml"

	"github.com/olive-io/bpmn/v2/pkg/data"
	"github.com/olive-io/bpmn/v2/pkg/errors"
	"github.com/olive-io/bpmn/v2/pkg/expression"
)

// XPath language engine
//
// Implementation details and limitations as per https://github.com/antchfx/xpath
type XPath struct {
	ctx               context.Context
	itemAwareLocators map[string]data.IItemAwareLocator
}

func (engine *XPath) SetItemAwareLocator(name string, itemAwareLocator data.IItemAwareLocator) {
	engine.itemAwareLocators[name] = itemAwareLocator
}

func Make(ctx context.Context) XPath {
	return XPath{ctx: ctx}
}

func New(ctx context.Context) *XPath {
	engine := Make(ctx)
	return &engine
}

func (engine *XPath) CompileExpression(source string) (result expression.ICompiledExpression, err error) {
	compiled, err := grammar.Build(source)
	if err == nil {
		result = &compiled
	}
	return
}

func (engine *XPath) EvaluateExpression(e expression.ICompiledExpression, datum interface{}) (result expression.IResult, err error) {
	if expr, ok := e.(*grammar.Grammar); ok {
		// Here, in order to save some prototyping type,
		// instead of implementing `parser.Parser` for `interface{}`,
		// we use it over `interface{}` serialized as XML.
		// This is not very efficient but does the job for now.
		// Eventually, a direct implementation of `parser.Parser`
		// over `interface{}` should be developed to optimize this path.

		var serialized []byte
		serialized, err = anyxml.Xml(datum)
		if err != nil {
			result = nil
			return
		}
		p := parser.ReadXml(bytes.NewBuffer(serialized))

		contextSettings := func(c *exec.ContextSettings) {
			if engine.itemAwareLocators != nil {
				c.FunctionLibrary[exec.XmlName{Local: data.LocatorObject}] = engine.getDataObject()
			}
		}

		var cursor store.Cursor
		cursor, err = store.CreateInMemory(p)
		if err != nil {
			return
		}
		var res exec.Result
		res, err = exec.Exec(cursor, expr, contextSettings)
		if err != nil {
			return
		}
		switch r := res.(type) {
		case exec.String:
			result = r.String()
		case exec.Bool:
			result = r.Bool()
		case exec.Number:
			result = r.Number()
		case exec.NodeSet:
			result = r
		}
	} else {
		err = errors.InvalidArgumentError{
			Expected: "CompiledExpression to be *github.com/ChrisTrenkamp/xsel/grammar#Grammar",
			Actual:   reflect.TypeOf(e),
		}
	}
	return
}

var asXMLType = reflect.TypeOf(new(data.IAsXML)).Elem()

func (engine *XPath) getDataObject() func(context exec.Context, args ...exec.Result) (exec.Result, error) {
	return func(context exec.Context, args ...exec.Result) (exec.Result, error) {
		var name string
		switch len(args) {
		case 0:
			return nil, errors.InvalidArgumentError{Expected: "at least one argument", Actual: "none"}
		case 2:
			return nil, errors.NotSupportedError{
				What:   "two-argument getDataObject",
				Reason: "doesn't support sub-processes yet",
			}
		case 1:
			name = args[0].String()
		}
		var itemAware data.IItemAware
		found := false
		for _, locator := range engine.itemAwareLocators {
			itemAware, found = locator.FindItemAwareByName(name)
			if !found {
				break
			}
		}

		if !found {
			return exec.NodeSet{}, nil
		}
		item := itemAware.Get()
		if item == nil {
			return nil, nil
		}
		switch tt := item.(type) {
		case string:
			return exec.String(tt), nil
		case float64:
			return exec.Number(tt), nil
		case bool:
			return exec.Bool(tt), nil
		default:
			// Until we have own data type to represent XML nodes, we'll piggy-back
			// on xsel's parser and datum.AsXML interface. This is not very efficient,
			// but should do for now
			if reflect.TypeOf(item).Implements(asXMLType) {
				p := parser.ReadXml(bytes.NewReader(item.(data.IAsXML).AsXML()))
				cursor, err := store.CreateInMemory(p)
				if err != nil {
					return nil, err
				}
				return exec.NodeSet{cursor}, nil
			} else {
				return nil, errors.InvalidArgumentError{
					Expected: "XML-serializable value (string, float64 or Node)",
					Actual:   item,
				}
			}
		}
	}
}

func init() {
	expression.RegisterEngine("http://www.w3.org/1999/XPath", func(ctx context.Context) expression.IEngine {
		return New(ctx)
	})
}
