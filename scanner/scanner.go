package scanner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/mkeesey/craftinginterpreters/token"
)

type Scanner struct {
	reader     *bufio.Reader
	tokens     []token.Token
	currLexeme bytes.Buffer

	line int
}

func NewScanner(reader io.Reader) *Scanner {
	read := bufio.NewReader(reader)
	buf := bytes.Buffer{}
	return &Scanner{reader: read, currLexeme: buf}
}

func (s *Scanner) scanTokens() ([]token.Token, error) {
	var allErrs error = nil
	for {
		s.currLexeme.Reset()
		err := s.scanToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			allErrs = errors.Join(allErrs, err)
		}
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, 0))
	return s.tokens, allErrs
}

func (s *Scanner) scanToken() error {
	rune, _, err := s.reader.ReadRune()
	if err != nil {
		return err
	}

	s.currLexeme.WriteRune(rune)
	switch rune {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	default:
		return fmt.Errorf("[%d] Error %s: %s", s.line, "", "Unexpected character")
	}
	return nil
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *Scanner) addTokenLiteral(tokenType token.TokenType, literal interface{}) {
	s.tokens = append(s.tokens, token.NewToken(tokenType, s.currLexeme.String(), literal, s.line))
}
