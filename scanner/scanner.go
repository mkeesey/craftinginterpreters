package scanner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

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
	var allErrs []error
	for {
		s.currLexeme.Reset()
		err := s.scanToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			allErrs = append(allErrs, err)
		}
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens, errors.Join(allErrs...)
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
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}
	case '/':
		if s.match('/') {
			for {
				bytes, err := s.reader.Peek(1)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					return fmt.Errorf("[%d] Error %s: %w", s.line, "err consuming comment", err)
				}
				if bytes[0] == '\n' {
					break
				}
				s.reader.Discard(1)
			}
		} else {
			s.addToken(token.SLASH)
		}
	case ' ', '\r', '\t':
		// ignore whitespace
	case '\n':
		s.line++
	case '"':
		err = s.stringToken()
	default:
		return fmt.Errorf("[%d] Error %s: %s", s.line, "", "Unexpected character")
	}
	return err
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *Scanner) addTokenLiteral(tokenType token.TokenType, literal interface{}) {
	s.tokens = append(s.tokens, token.NewToken(tokenType, s.currLexeme.String(), literal, s.line))
}

func (s *Scanner) match(expected rune) bool {
	runeLength := utf8.RuneLen(expected)
	bytes, err := s.reader.Peek(runeLength)
	if err != nil {
		return false
	}
	seen, _ := utf8.DecodeRune(bytes)
	if seen != expected {
		return false
	}
	s.currLexeme.WriteRune(seen)
	s.reader.Discard(runeLength)
	return true
}

func (s *Scanner) stringToken() error {
	s.currLexeme.Reset() // remove leading quote
	for {
		bytes, err := s.reader.Peek(1)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return fmt.Errorf("[%d] Error %s: %s", s.line, "Unterminated string", s.currLexeme.String())
			}
			return fmt.Errorf("[%d] Error %s: %w", s.line, "err consuming string", err)
		}
		if bytes[0] == '"' {
			s.reader.Discard(1) // skip closing quote
			break
		}

		if bytes[0] == '\n' {
			s.line++
		}
		rune, _, err := s.reader.ReadRune()
		if err != nil {
			return fmt.Errorf("[%d] Error %s: %w", s.line, "err consuming string", err)
		}
		s.currLexeme.WriteRune(rune)
	}

	s.addTokenLiteral(token.STRING, s.currLexeme.String())

	return nil
}
