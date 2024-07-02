package parser

import (
	"errors"

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
	allErrs := []error{}
	for {
		// TODO - peek is not in the book, but it's necessary to check for EOF
		if p.isAtEnd() || p.peek().Type == token.EOF {
			break
		}
		stmt, err := p.declaration()
		if err != nil {
			allErrs = append(allErrs, err)
		}
		stmts = append(stmts, stmt)
	}

	return stmts, errors.Join(allErrs...)
}

func (p *Parser) declaration() (ast.Stmt, error) {
	var stmt ast.Stmt
	var err error
	if p.match(token.VAR) {
		stmt, err = p.varDeclaration()
	} else {
		stmt, err = p.statement()
	}

	if err != nil {
		p.synchronize()
	}
	return stmt, err
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return &ast.StmtVar{Name: name, Initializer: initializer}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	} else if p.match(token.LEFT_BRACE) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.Block{Statements: stmts}, nil
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

func (p *Parser) block() ([]ast.Stmt, error) {
	stmts := []ast.Stmt{}
	for {
		if p.check(token.RIGHT_BRACE) || p.isAtEnd() {
			break
		}
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return stmts, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.equality()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if exprVar, ok := expr.(*ast.ExprVar); ok {
			name := exprVar.Name
			return &ast.Assign{Name: name, Value: value}, nil
		}

		return nil, failure.TokenError(equals, "Invalid assignment target.")
	}

	return expr, err
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
	} else if p.match(token.IDENTIFIER) {
		return &ast.ExprVar{Name: p.previous()}, nil
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
			return
		default:
			p.advance()
		}
	}
}
