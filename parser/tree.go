package parser

import (
	"fmt"
	"github.com/keatonmartin/golox/token"
)

type RuntimeError struct {
	Message string
	Tok     token.Token
}

type Expr interface {
	String() string
	Interpret() interface{}
}

type Binary struct {
	Left, Right Expr
	Operator    token.Token
}

func (b Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Operator.Lexeme, b.Right.String())
}

func (b Binary) Interpret() interface{} {
	left := b.Left.Interpret()
	right := b.Right.Interpret()

	switch b.Operator.TokenType {
	case token.MINUS:
		assertNumberOperands(b.Operator, left, right)
		return left.(float64) - right.(float64)
	case token.PLUS:
		_, ok1 := left.(float64)
		_, ok2 := right.(float64)
		if ok1 && ok2 {
			return left.(float64) + right.(float64)
		}
		_, ok1 = left.(string)
		_, ok2 = right.(string)
		if ok1 && ok2 {
			return left.(string) + right.(string)
		}
		panic(RuntimeError{"Operands must be two numbers or two strings", b.Operator})
	case token.STAR:
		assertNumberOperands(b.Operator, left, right)
		return left.(float64) * right.(float64)
	case token.SLASH:
		assertNumberOperands(b.Operator, left, right)
		return left.(float64) / right.(float64)
	case token.GREATER:
		assertNumberOperands(b.Operator, left, right)
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		assertNumberOperands(b.Operator, left, right)
		return left.(float64) >= right.(float64)
	case token.LESS:
		assertNumberOperands(b.Operator, left, right)
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		assertNumberOperands(b.Operator, left, right)
		return left.(float64) <= right.(float64)
	case token.EQUAL_EQUAL:
		return left == right
	case token.BANG_EQUAL:
		return !(left == right)
	}
	return nil // unreachable
}

type Grouping struct {
	Expr Expr
}

func (g Grouping) String() string {
	return fmt.Sprintf("(%s)", g.Expr.String())
}

func (g Grouping) Interpret() interface{} {
	return g.Expr.Interpret()
}

type Literal struct {
	Value interface{}
}

func (l Literal) String() string {
	return fmt.Sprintf("%v", l.Value)
}

func (l Literal) Interpret() interface{} {
	return l.Value
}

type Unary struct {
	Right    Expr
	Operator token.Token
}

func (u Unary) String() string {
	return fmt.Sprintf("(%s%s)", u.Operator.Lexeme, u.Right.String())
}

func (u Unary) Interpret() interface{} {
	right := u.Right.Interpret()
	switch u.Operator.TokenType {
	case token.MINUS:
		assertNumberOperands(u.Operator, right)
		return -1 * (right.(float64))
	case token.BANG:
		return !isTruthy(u.Right)
	}
	return nil
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	b, ok := v.(bool)
	if ok {
		return b
	}
	return true
}

// check that all operands are float64. panic with runtime error if not.
func assertNumberOperands(operator token.Token, operands ...interface{}) {
	for _, op := range operands {
		_, ok := op.(float64)
		if !ok {
			panic(RuntimeError{"Operands must be numbers.", operator})
		}
	}
}
