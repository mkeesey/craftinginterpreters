package token

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{Type: tokenType, Lexeme: lexeme, Literal: literal, Line: line}
}

func (t *Token) String() string {
	return fmt.Sprintf("{%v %s %v}", t.Type, t.Lexeme, t.Literal)
}
