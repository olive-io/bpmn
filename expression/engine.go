package expression

import "github.com/olive-io/bpmn/data"

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
