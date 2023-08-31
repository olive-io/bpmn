package expr

import (
	"context"
	"reflect"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/olive-io/bpmn/data"
	"github.com/olive-io/bpmn/errors"
	"github.com/olive-io/bpmn/expression"
)

// Expr language engine
//
// https://github.com/antonmedv/expr
type Expr struct {
	env              map[string]interface{}
	itemAwareLocator data.IItemAwareLocator
}

func (engine *Expr) SetItemAwareLocator(itemAwareLocator data.IItemAwareLocator) {
	engine.itemAwareLocator = itemAwareLocator
}

func New(ctx context.Context) *Expr {
	engine := &Expr{}
	engine.env = map[string]interface{}{
		"getDataObject": func(args ...string) data.IItem {
			var name string
			if len(args) == 1 {
				name = args[0]
			}
			itemAware, found := engine.itemAwareLocator.FindItemAwareByName(name)
			if !found {
				return nil
			}
			ch := itemAware.Get(ctx)
			if ch == nil {
				return nil
			}
			select {
			case <-ctx.Done():
				return nil
			case item := <-ch:
				return item
			}
		},
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
			Expected: "CompiledExpression to be *github.com/antonmedv/expr/vm#Program",
			Actual:   reflect.TypeOf(e),
		}
	}
	return
}

func init() {
	expression.RegisterEngine("https://github.com/antonmedv/expr", func(ctx context.Context) expression.IEngine {
		return New(ctx)
	})
}
