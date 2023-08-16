package interpreter

import "github.com/keatonmartin/golox/parser"

type Interpreter struct {
	Expr parser.Expr // top level expression to evaluate
	Errs []parser.RuntimeError
}

func NewInterpreter(expr parser.Expr) Interpreter {
	return Interpreter{expr, []parser.RuntimeError{}}
}

func (i *Interpreter) Interpret() interface{} {
	defer func() {
		err := recover()
		if err != nil {
			i.Errs = append(i.Errs, err.(parser.RuntimeError))
		}
	}()
	res := i.Expr.Interpret()
	return res
}
