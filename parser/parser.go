package parser

import (
	"github.com/keatonmartin/golox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
	Errs    []parseError
}

type parseError struct {
	Token   token.Token
	Message string
}

func NewParser(tokens []token.Token) Parser {
	return Parser{tokens, 0, []parseError{}}
}

func (p *Parser) Parse() Expr {
	defer func() {
		recover()
	}()
	exp := p.expression()
	return exp
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{expr, right, operator}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Binary{expr, right, operator}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{expr, right, operator}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{expr, right, operator}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary{right, operator}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(token.FALSE) {
		return Literal{false}
	} else if p.match(token.TRUE) {
		return Literal{true}
	} else if p.match(token.NIL) {
		return Literal{nil}
	}

	if p.match(token.NUMBER, token.STRING) {
		return Literal{p.previous().Literal}
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{expr}
	}
	p.Errs = append(p.Errs, parseError{p.peek(), "Expect expression."})
	panic("Expect expression") // maybe not the cleanest solution
}

func (p *Parser) consume(t token.Type, message string) token.Token {
	if p.check(t) {
		return p.advance()
	}
	p.Errs = append(p.Errs, parseError{p.peek(), message})
	panic(message)
}

// match will determine if the current token is any type t in types
func (p *Parser) match(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

// check will check if the current token is of type t
func (p *Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

// isAtEnd returns if all tokens have been parsed
func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

// peek returns the next token
func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

// advance "consumes" the current token
func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// previous returns the previous token
func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().TokenType == token.SEMICOLON {
			return
		}
		switch p.peek().TokenType {
		case token.CLASS, token.FOR, token.FUN, token.IF, token.PRINT, token.RETURN, token.VAR, token.WHILE:
			return
		}
	}
}
