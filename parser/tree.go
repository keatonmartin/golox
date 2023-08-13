package parser

import (
	"fmt"
	"github.com/keatonmartin/golox/token"
)

// Each struct in this file is a node in an AST

type Expr interface {
	String() string
}

type Binary struct {
	Left, Right Expr
	Operator    token.Token
}

func (b Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Operator.Lexeme, b.Right.String())
}

type Grouping struct {
	Expr Expr
}

func (g Grouping) String() string {
	return fmt.Sprintf("(%s)", g.Expr.String())
}

type Literal struct {
	Value interface{}
}

func (l Literal) String() string {
	return fmt.Sprintf("%v", l.Value)
}

type Unary struct {
	Right    Expr
	Operator token.Token
}

func (u Unary) String() string {
	return fmt.Sprintf("(%s%s)", u.Operator.Lexeme, u.Right.String())
}

type Ternary struct {
	Conditional, Left, Right Expr
}

func (t Ternary) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", t.Conditional.String(), t.Left.String(), t.Right.String())
}
