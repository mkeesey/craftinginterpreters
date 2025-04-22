package bytecode

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR

	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL

	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FOR
	TOKEN_FUN
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_ERROR
	TOKEN_EOF
)

func (t TokenType) String() string {
	switch t {
	// Single-character tokens
	case TOKEN_LEFT_PAREN:
		return "TOKEN_LEFT_PAREN"
	case TOKEN_RIGHT_PAREN:
		return "TOKEN_RIGHT_PAREN"
	case TOKEN_LEFT_BRACE:
		return "TOKEN_LEFT_BRACE"
	case TOKEN_RIGHT_BRACE:
		return "TOKEN_RIGHT_BRACE"
	case TOKEN_COMMA:
		return "TOKEN_COMMA"
	case TOKEN_DOT:
		return "TOKEN_DOT"
	case TOKEN_MINUS:
		return "TOKEN_MINUS"
	case TOKEN_PLUS:
		return "TOKEN_PLUS"
	case TOKEN_SEMICOLON:
		return "TOKEN_SEMICOLON"
	case TOKEN_SLASH:
		return "TOKEN_SLASH"
	case TOKEN_STAR:
		return "TOKEN_STAR"

	// One or two character tokens
	case TOKEN_BANG:
		return "TOKEN_BANG"
	case TOKEN_BANG_EQUAL:
		return "TOKEN_BANG_EQUAL"
	case TOKEN_EQUAL:
		return "TOKEN_EQUAL"
	case TOKEN_EQUAL_EQUAL:
		return "TOKEN_EQUAL_EQUAL"
	case TOKEN_GREATER:
		return "TOKEN_GREATER"
	case TOKEN_GREATER_EQUAL:
		return "TOKEN_GREATER_EQUAL"
	case TOKEN_LESS:
		return "TOKEN_LESS"
	case TOKEN_LESS_EQUAL:
		return "TOKEN_LESS_EQUAL"

	// Literals
	case TOKEN_IDENTIFIER:
		return "TOKEN_IDENTIFIER"
	case TOKEN_STRING:
		return "TOKEN_STRING"
	case TOKEN_NUMBER:
		return "TOKEN_NUMBER"

	// Keywords
	case TOKEN_AND:
		return "TOKEN_AND"
	case TOKEN_CLASS:
		return "TOKEN_CLASS"
	case TOKEN_ELSE:
		return "TOKEN_ELSE"
	case TOKEN_FALSE:
		return "TOKEN_FALSE"
	case TOKEN_FOR:
		return "TOKEN_FOR"
	case TOKEN_FUN:
		return "TOKEN_FUN"
	case TOKEN_IF:
		return "TOKEN_IF"
	case TOKEN_NIL:
		return "TOKEN_NIL"
	case TOKEN_OR:
		return "TOKEN_OR"
	case TOKEN_PRINT:
		return "TOKEN_PRINT"
	case TOKEN_RETURN:
		return "TOKEN_RETURN"
	case TOKEN_SUPER:
		return "TOKEN_SUPER"
	case TOKEN_THIS:
		return "TOKEN_THIS"
	case TOKEN_TRUE:
		return "TOKEN_TRUE"
	case TOKEN_VAR:
		return "TOKEN_VAR"
	case TOKEN_WHILE:
		return "TOKEN_WHILE"

	// Special tokens
	case TOKEN_ERROR:
		return "TOKEN_ERROR"
	case TOKEN_EOF:
		return "TOKEN_EOF"
	default:
		return fmt.Sprintf("TOKEN_UNKNOWN(%d)", t)
	}
}

type token struct {
	tokenType TokenType
	lexeme    string
	line      int
}

type scanner struct {
	line       int
	startIdx   int
	currentIdx int
	source     string
}

func newScanner(source string) *scanner {
	return &scanner{
		line:       1,
		startIdx:   0,
		currentIdx: 0,
		source:     source,
	}
}

func (s *scanner) scanToken() *token {
	s.skipWhitespace()
	s.startIdx = s.currentIdx
	if s.isAtEnd() {
		return s.makeToken(TOKEN_EOF)
	}

	r := s.advance()
	if unicode.IsLetter(r) || r == '_' {
		return s.identifier()
	}
	if unicode.IsDigit(r) {
		return s.number()
	}

	switch r {
	case '(':
		return s.makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return s.makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return s.makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return s.makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return s.makeToken(TOKEN_SEMICOLON)
	case ',':
		return s.makeToken(TOKEN_COMMA)
	case '.':
		return s.makeToken(TOKEN_DOT)
	case '-':
		return s.makeToken(TOKEN_MINUS)
	case '+':
		return s.makeToken(TOKEN_PLUS)
	case '/':
		return s.makeToken(TOKEN_SLASH)
	case '*':
		return s.makeToken(TOKEN_STAR)
	case '!':
		if s.match('=') {
			return s.makeToken(TOKEN_BANG_EQUAL)
		}
		return s.makeToken(TOKEN_BANG)
	case '=':
		if s.match('=') {
			return s.makeToken(TOKEN_EQUAL_EQUAL)
		}
		return s.makeToken(TOKEN_EQUAL)
	case '<':
		if s.match('=') {
			return s.makeToken(TOKEN_LESS_EQUAL)
		}
		return s.makeToken(TOKEN_LESS)
	case '>':
		if s.match('=') {
			return s.makeToken(TOKEN_GREATER_EQUAL)
		}
	case '"':
		return s.string()
	}

	return s.errorToken("Unrecognized character.")
}

func (s *scanner) skipWhitespace() {
	for {
		r := s.peek()
		switch r {
		case ' ', '\r', '\t':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		case '/':
			if s.peekNext() == '/' {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (s *scanner) string() *token {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken("Unterminated string.")
	}

	s.advance() // The closing '"'.

	return s.makeToken(TOKEN_STRING)
}

func (s *scanner) identifier() *token {
	for isAlpha(s.peek()) || unicode.IsDigit(s.peek()) {
		s.advance()
	}
	return s.makeToken(s.identifierType())
}

func (s *scanner) identifierType() TokenType {
	switch s.source[s.startIdx] {
	case 'a':
		return s.checkKeyword(s.startIdx, 3, "and", TOKEN_AND)
	case 'c':
		return s.checkKeyword(s.startIdx, 5, "class", TOKEN_CLASS)
	case 'e':
		return s.checkKeyword(s.startIdx, 4, "else", TOKEN_ELSE)
	case 'i':
		return s.checkKeyword(s.startIdx, 2, "if", TOKEN_IF)
	case 'n':
		return s.checkKeyword(s.startIdx, 3, "nil", TOKEN_NIL)
	case 'o':
		return s.checkKeyword(s.startIdx, 2, "or", TOKEN_OR)
	case 'p':
		return s.checkKeyword(s.startIdx, 5, "print", TOKEN_PRINT)
	case 'r':
		return s.checkKeyword(s.startIdx, 6, "return", TOKEN_RETURN)
	case 's':
		return s.checkKeyword(s.startIdx, 5, "super", TOKEN_SUPER)
	case 'v':
		return s.checkKeyword(s.startIdx, 3, "var", TOKEN_VAR)
	case 'w':
		return s.checkKeyword(s.startIdx, 5, "while", TOKEN_WHILE)
	case 'f':
		if s.currentIdx-s.startIdx > 1 {
			switch s.source[s.startIdx+1] {
			case 'a':
				return s.checkKeyword(s.startIdx, 4, "false", TOKEN_FALSE)
			case 'o':
				return s.checkKeyword(s.startIdx, 3, "for", TOKEN_FOR)
			case 'u':
				return s.checkKeyword(s.startIdx, 3, "fun", TOKEN_FUN)
			}
		}
	case 't':
		if s.currentIdx-s.startIdx > 1 {
			switch s.source[s.startIdx+1] {
			case 'h':
				return s.checkKeyword(s.startIdx, 4, "this", TOKEN_THIS)
			case 'r':
				return s.checkKeyword(s.startIdx, 4, "true", TOKEN_TRUE)
			}
		}
	}
	return TOKEN_IDENTIFIER
}

func (s *scanner) checkKeyword(start int, length int, rest string, type_ TokenType) TokenType {
	if s.currentIdx == start+length && s.source[start:start+length] == rest {
		return type_
	}
	return TOKEN_IDENTIFIER
}

func (s *scanner) number() *token {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()
		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	return s.makeToken(TOKEN_NUMBER)
}

func (s *scanner) peek() rune {
	r, _ := utf8.DecodeRuneInString(s.source[s.currentIdx:])
	return r
}

func (s *scanner) peekNext() rune {
	if s.isAtEnd() {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(s.source[s.currentIdx+1:])
	return r
}

func (s *scanner) advance() rune {
	r, size := utf8.DecodeRuneInString(s.source[s.currentIdx:])
	s.currentIdx += size

	return r
}

func (s *scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.currentIdx] != byte(expected) {
		return false
	}
	s.currentIdx++
	return true
}

func (s *scanner) makeToken(tokenType TokenType) *token {
	lexeme := s.source[s.startIdx:s.currentIdx]
	return &token{
		tokenType: tokenType,
		lexeme:    lexeme,
		line:      s.line,
	}
}

func (s *scanner) errorToken(message string) *token {
	return &token{
		tokenType: TOKEN_ERROR,
		lexeme:    message,
		line:      s.line,
	}
}

func (s *scanner) isAtEnd() bool {
	return s.currentIdx >= len(s.source)
}

func compile(source string) {
	scanner := newScanner(source)
	line := -1
	for {
		token := scanner.scanToken()
		if token.line != line {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Print("   | ")
		}

		fmt.Printf("%s '%s'\n", token.tokenType, token.lexeme)

		if token.tokenType == TOKEN_EOF {
			break
		}
	}
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}
