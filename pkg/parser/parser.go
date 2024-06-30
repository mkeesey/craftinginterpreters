package parser

import (
	"github.com/mkeesey/craftinginterpreters/pkg/ast"
	"github.com/mkeesey/craftinginterpreters/pkg/failure"
	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	stmts := []ast.Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err // TODO, should accumulate probably
		}
		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &ast.Print{Expression: expr}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &ast.Expression{Expression: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()

	for {
		if !p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
			break
		}
		operator := p.previous()
		var right ast.Expr
		right, err = p.comparison()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, err
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()

	for {
		if !p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
			break
		}
		operator := p.previous()
		var right ast.Expr
		right, err = p.term()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, err
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()

	for {
		if !p.match(token.MINUS, token.PLUS) {
			break
		}
		operator := p.previous()
		var right ast.Expr
		right, err = p.factor()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, err
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()

	for {
		if !p.match(token.SLASH, token.STAR) {
			break
		}
		operator := p.previous()
		var right ast.Expr
		right, err = p.unary()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, err
}

func (p *Parser) unary() (ast.Expr, error) {
	for {
		if !p.match(token.BANG, token.MINUS) {
			break
		}
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{Operator: operator, Right: right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return &ast.Literal{Value: false}, nil
	} else if p.match(token.TRUE) {
		return &ast.Literal{Value: true}, nil
	} else if p.match(token.NIL) {
		return &ast.Literal{Value: nil}, nil
	} else if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{Value: p.previous().Literal}, nil
	} else if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &ast.Grouping{Expression: expr}, nil
	}

	return nil, failure.TokenError(p.peek(), "Expect expression.")
}

func (p *Parser) consume(t token.TokenType, message string) (*token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return nil, failure.TokenError(p.peek(), message)
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

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
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

func (p *Parser) synchronize() {
	p.advance()

	for {
		if p.isAtEnd() || p.previous().Type == token.SEMICOLON {
			break
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			break
		default:
			p.advance()
		}
	}
}
