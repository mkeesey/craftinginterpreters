package bytecode

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

var debugPrintCode = false

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

	TOKEN_MAX
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

type precedence int

const (
	PREC_NONE       precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

type parsefn func()

type parserule struct {
	prefix     parsefn
	infix      parsefn
	precedence precedence
}

type parser struct {
	hadError       bool
	panicMode      bool
	rules          []parserule
	scanner        *scanner
	current        *token
	previous       *token
	compilingChunk *Chunk
}

func newParser(scanner *scanner, chunk *Chunk) *parser {
	parser := &parser{
		scanner:        scanner,
		compilingChunk: chunk,
	}

	rules := make([]parserule, TOKEN_MAX)

	// Initialize all rules to nil functions and PREC_NONE by default
	for i := range rules {
		rules[i] = parserule{nil, nil, PREC_NONE}
	}

	// Set specific rules
	rules[TOKEN_LEFT_PAREN] = parserule{parser.grouping, nil, PREC_NONE}
	rules[TOKEN_MINUS] = parserule{parser.unary, parser.binary, PREC_TERM}
	rules[TOKEN_PLUS] = parserule{nil, parser.binary, PREC_TERM}
	rules[TOKEN_SLASH] = parserule{nil, parser.binary, PREC_FACTOR}
	rules[TOKEN_STAR] = parserule{nil, parser.binary, PREC_FACTOR}
	rules[TOKEN_NUMBER] = parserule{parser.number, nil, PREC_NONE}
	rules[TOKEN_FALSE] = parserule{parser.literal, nil, PREC_NONE}
	rules[TOKEN_TRUE] = parserule{parser.literal, nil, PREC_NONE}
	rules[TOKEN_NIL] = parserule{parser.literal, nil, PREC_NONE}
	rules[TOKEN_BANG] = parserule{parser.unary, nil, PREC_NONE}
	rules[TOKEN_BANG_EQUAL] = parserule{nil, parser.binary, PREC_EQUALITY}
	rules[TOKEN_EQUAL_EQUAL] = parserule{nil, parser.binary, PREC_EQUALITY}
	rules[TOKEN_GREATER] = parserule{nil, parser.binary, PREC_COMPARISON}
	rules[TOKEN_GREATER_EQUAL] = parserule{nil, parser.binary, PREC_COMPARISON}
	rules[TOKEN_LESS] = parserule{nil, parser.binary, PREC_COMPARISON}
	rules[TOKEN_LESS_EQUAL] = parserule{nil, parser.binary, PREC_COMPARISON}

	parser.rules = rules
	return parser
}

func (p *parser) advance() {
	p.previous = p.current

	for {
		p.current = p.scanner.scanToken()
		if p.current.tokenType != TOKEN_ERROR {
			break
		}

		p.errorAtCurrent(p.current.lexeme)
	}
}

func (p *parser) consume(tokenType TokenType, msg string) {
	if p.current.tokenType == tokenType {
		p.advance()
		return
	}

	p.errorAtCurrent(msg)
}

func (p *parser) end() {
	p.emitByte(byte(OP_RETURN))

	if debugPrintCode {
		if !p.hadError {
			DisassembleChunk(p.compilingChunk, "code")
		}
	}
}

func (p *parser) number() {
	val, err := strconv.ParseFloat(p.previous.lexeme, 64)
	if err != nil {
		panic(err)
	}
	p.emitBytes(byte(OP_CONSTANT), byte(p.makeConstant(NumberValue(val))))
}

func (p *parser) literal() {
	switch p.previous.tokenType {
	case TOKEN_FALSE:
		p.emitByte(byte(OP_FALSE))
	case TOKEN_TRUE:
		p.emitByte(byte(OP_TRUE))
	case TOKEN_NIL:
		p.emitByte(byte(OP_NIL))
	default:
		return // unreachable
	}
}

func (p *parser) grouping() {
	p.expression()
	p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression")
}

func (p *parser) unary() {
	opType := p.previous.tokenType

	p.expression()

	switch opType {
	case TOKEN_MINUS:
		p.emitByte(byte(OP_NEGATE))
	case TOKEN_BANG:
		p.emitByte(byte(OP_NOT))
	default:
		return // unreachable
	}
}

func (p *parser) binary() {
	opType := p.previous.tokenType
	rule := p.getRule(opType)
	p.parsePrecedence(rule.precedence + 1)

	switch opType {
	case TOKEN_BANG_EQUAL:
		p.emitBytes(byte(OP_EQUAL), byte(OP_NOT))
	case TOKEN_EQUAL_EQUAL:
		p.emitByte(byte(OP_EQUAL))
	case TOKEN_GREATER:
		p.emitByte(byte(OP_GREATER))
	case TOKEN_GREATER_EQUAL:
		p.emitBytes(byte(OP_LESS), byte(OP_NOT))
	case TOKEN_LESS:
		p.emitByte(byte(OP_LESS))
	case TOKEN_LESS_EQUAL:
		p.emitBytes(byte(OP_GREATER), byte(OP_NOT))
	case TOKEN_PLUS:
		p.emitByte(byte(OP_ADD))
	case TOKEN_MINUS:
		p.emitByte(byte(OP_SUBTRACT))
	case TOKEN_STAR:
		p.emitByte(byte(OP_MULTIPLY))
	case TOKEN_SLASH:
		p.emitByte(byte(OP_DIVIDE))
	}
}

func (p *parser) parsePrecedence(precedence precedence) {
	p.advance()
	prefixRule := p.getRule(p.previous.tokenType).prefix
	if prefixRule == nil {
		p.error("Expect expression.")
		return
	}

	prefixRule()

	for precedence <= p.getRule(p.current.tokenType).precedence {
		p.advance()
		infixRule := p.getRule(p.previous.tokenType).infix
		infixRule()
	}
}

func (p *parser) getRule(tokenType TokenType) parserule {
	return p.rules[tokenType]
}

func (p *parser) expression() {
	p.parsePrecedence(PREC_ASSIGNMENT)
}

func (p *parser) makeConstant(val Value) uint8 {
	constant := p.compilingChunk.WriteConstant(val)
	if constant >= 255 { // TODO max uint8 size
		p.error("Too many constants in one chunk.")
		return 0
	}

	return constant
}

func (p *parser) emitByte(val byte) {
	p.compilingChunk.Write(val, p.previous.line)
}

func (p *parser) emitBytes(valOne byte, valTwo byte) {
	p.emitByte(valOne)
	p.emitByte(valTwo)
}

func (p *parser) errorAtCurrent(msg string) {
	p.errorAt(p.current, msg)
}

func (p *parser) error(msg string) {
	p.errorAt(p.previous, msg)
}

func (p *parser) errorAt(tok *token, msg string) {
	if p.panicMode {
		return
	}

	p.panicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", tok.line)

	if tok.tokenType == TOKEN_EOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if tok.tokenType == TOKEN_ERROR {
		// Nothing
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", tok.lexeme)
	}

	fmt.Fprintf(os.Stderr, ": %s\n", msg)
	p.hadError = true
}

func compile(source string, chunk *Chunk) error {
	scanner := newScanner(source)
	parser := newParser(scanner, chunk)
	parser.advance()
	parser.expression()
	parser.consume(TOKEN_EOF, "Expect end of expression")
	parser.end()

	if parser.hadError {
		return fmt.Errorf("Parsing error")
	}
	return nil
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}
