// Copyright 2023 The bpmn Authors
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
	"reflect"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/expression"
)

// Expr language engine
//
// https://github.com/expr-lang/expr
type Expr struct {
	ctx               context.Context
	itemAwareLocators map[string]data.IItemAwareLocator
	env               map[string]interface{}
}

func (engine *Expr) SetItemAwareLocator(name string, itemAwareLocator data.IItemAwareLocator) {
	engine.itemAwareLocators[name] = itemAwareLocator
}

func (engine *Expr) fetchItem(pool string) func(args ...string) data.IItem {
	return func(args ...string) data.IItem {
		var name string
		if len(args) == 1 {
			name = args[0]
		}
		locator, ok := engine.itemAwareLocators[pool]
		if !ok {
			return nil
		}

		itemAware, found := locator.FindItemAwareByName(name)
		if !found {
			return nil
		}
		item := itemAware.Get()
		return item
	}
}

func New(ctx context.Context) *Expr {
	engine := &Expr{
		ctx:               ctx,
		itemAwareLocators: map[string]data.IItemAwareLocator{},
	}
	engine.env = map[string]interface{}{
		data.LocatorObject: engine.fetchItem(data.LocatorObject),
		data.LocatorHeader: engine.fetchItem(data.LocatorHeader),
	}
	return engine
}

func (engine *Expr) CompileExpression(source string) (result expression.ICompiledExpression, err error) {
	result, err = expr.Compile(source, expr.Env(engine.env), expr.AllowUndefinedVariables())
	return
}

func (engine *Expr) EvaluateExpression(e expression.ICompiledExpression,
	data interface{},
) (result expression.IResult, err error) {
	actualData := data
	if data == nil {
		actualData = engine.env
	} else if e, ok := data.(map[string]interface{}); ok {
		env := engine.env
		for k, v := range e {
			env[k] = v
		}
		actualData = env
	}

	if exp, ok := e.(*vm.Program); ok {
		result, err = expr.Run(exp, actualData)
	} else {
		err = errors.InvalidArgumentError{
			Expected: "CompiledExpression to be *github.com/expr-lang/expr/vm#Program",
			Actual:   reflect.TypeOf(e),
		}
	}
	return
}

func init() {
	expression.RegisterEngine("https://github.com/expr-lang/expr", func(ctx context.Context) expression.IEngine {
		return New(ctx)
	})
}
