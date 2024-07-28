package scanner

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mkeesey/craftinginterpreters/pkg/failure"
	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

type Scanner struct {
	reader     *bufio.Reader
	tokens     []*token.Token
	currLexeme strings.Builder

	line int

	reporter *failure.Reporter
}

func NewScanner(reader io.Reader, reporter *failure.Reporter) *Scanner {
	read := bufio.NewReader(reader)
	buf := strings.Builder{}
	return &Scanner{reader: read, currLexeme: buf, reporter: reporter, line: 1}
}

func (s *Scanner) ScanTokens() []*token.Token {
	for {
		s.currLexeme.Reset()
		err := s.scanToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				s.reporter.Panic(s.line, err)
			}
		}
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens
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
					s.reporter.Panic(s.line, err)
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
		if isNumber(rune) {
			s.numberToken()
		} else if isAlpha(rune) {
			err = s.identifierToken()
		} else {
			s.reporter.Error(s.line, "Unexpected character.")
		}
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
				s.reporter.Error(s.line, "Unterminated string.")
				return nil
			}
			s.reporter.Panic(s.line, err)
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
			s.reporter.Panic(s.line, err)
		}
		s.currLexeme.WriteRune(rune)
	}

	s.addTokenLiteral(token.STRING, s.currLexeme.String())

	return nil
}

func (s *Scanner) numberToken() error {
	err := s.consumeDigits()
	if err != nil {
		return err
	}

	bytes, err := s.reader.Peek(2)
	if err == nil {
		if bytes[0] == '.' && isByteNumber(bytes[1]) {
			s.reader.Discard(1)
			s.currLexeme.WriteRune('.')
			err = s.consumeDigits()
			if err != nil {
				return err
			}
		}
	}

	literal, err := strconv.ParseFloat(s.currLexeme.String(), 64)
	if err != nil {
		s.reporter.Panic(s.line, err)
	}

	s.addTokenLiteral(token.NUMBER, literal)
	return nil
}

func (s *Scanner) consumeDigits() error {
	for {
		bytes, err := s.reader.Peek(1)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			s.reporter.Panic(s.line, err)
		}
		if !isByteNumber(bytes[0]) {
			break
		}
		rune, _, err := s.reader.ReadRune()
		if err != nil {
			s.reporter.Panic(s.line, err)
		}
		s.currLexeme.WriteRune(rune)
	}
	return nil
}

func (s *Scanner) identifierToken() error {
	for {
		rune, _, err := s.reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return failure.Wrap(s.line, "err consuming identifier rune", err)
		}
		if !isAlphaNumeric(rune) {
			s.reader.UnreadRune()
			break
		}
		s.currLexeme.WriteRune(rune)
	}

	identifier := s.currLexeme.String()
	tokenType, ok := keywords[identifier]
	if ok {
		s.addToken(tokenType)
	} else {
		s.addToken(token.IDENTIFIER)
	}

	return nil
}

func isNumber(rune rune) bool {
	return rune >= '0' && rune <= '9'
}

func isByteNumber(val byte) bool {
	return val >= '0' && val <= '9'
}

func isAlpha(rune rune) bool {
	return unicode.IsLetter(rune) || rune == '_'
}

func isAlphaNumeric(rune rune) bool {
	return isAlpha(rune) || isNumber(rune)
}
