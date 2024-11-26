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

package expression

import (
	"github.com/olive-io/bpmn/v2/pkg/data"
)

type ICompiler interface {
	CompileExpression(source string) (ICompiledExpression, error)
}

type ICompiledExpression interface{}

type IEvaluator interface {
	EvaluateExpression(expr ICompiledExpression, data interface{}) (IResult, error)
}

type IResult interface{}

type IEngine interface {
	ICompiler
	IEvaluator
	SetItemAwareLocator(name string, itemAwareLocator data.IItemAwareLocator)
}
