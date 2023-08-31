package expression

import "github.com/olive-io/bpmn/data"

type Compiler interface {
	CompileExpression(source string) (CompiledExpression, error)
}

type CompiledExpression interface{}

type Evaluator interface {
	EvaluateExpression(expr CompiledExpression, data interface{}) (Result, error)
}

type Result interface{}

type Engine interface {
	Compiler
	Evaluator
	SetItemAwareLocator(itemAwareLocator data.ItemAwareLocator)
}
