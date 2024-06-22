package scanner

import "github.com/mkeesey/craftinginterpreters/token"

var (
	keywords = map[string]token.TokenType{
		"and":    token.AND,
		"class":  token.CLASS,
		"else":   token.ELSE,
		"false":  token.FALSE,
		"for":    token.FOR,
		"fun":    token.FUN,
		"if":     token.IF,
		"nil":    token.NIL,
		"or":     token.OR,
		"print":  token.PRINT,
		"return": token.RETURN,
		"super":  token.SUPER,
		"this":   token.THIS,
		"true":   token.TRUE,
		"var":    token.VAR,
		"while":  token.WHILE,
	}
)
