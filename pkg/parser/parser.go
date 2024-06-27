package parser

import (
	"fmt"

	"github.com/mkeesey/craftinginterpreters/pkg/ast"
	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for {
		if !p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
			break
		}
		operator := p.previous()
		right := p.comparison()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for {
		if !p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
			break
		}
		operator := p.previous()
		right := p.term()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for {
		if !p.match(token.MINUS, token.PLUS) {
			break
		}
		operator := p.previous()
		right := p.factor()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for {
		if !p.match(token.SLASH, token.STAR) {
			break
		}
		operator := p.previous()
		right := p.unary()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	for {
		if !p.match(token.BANG, token.MINUS) {
			break
		}
		operator := p.previous()
		right := p.unary()
		return &ast.Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return &ast.Literal{Value: false}
	} else if p.match(token.TRUE) {
		return &ast.Literal{Value: true}
	} else if p.match(token.NIL) {
		return &ast.Literal{Value: nil}
	} else if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{Value: p.previous().Literal}
	} else if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.Grouping{Expression: expr}
	}

	return nil //TODO
}

func (p *Parser) consume(t token.TokenType, message string) (*token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	// TODO error actually reported correctly
	return nil, fmt.Errorf(message)
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}
func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.tokens[p.current].Type == t
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}
