package bytecode

import (
	"strings"
	"testing"
)

func TestScanToken(t *testing.T) {
	// Test cases for different tokens
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"(", TOKEN_LEFT_PAREN},
		{")", TOKEN_RIGHT_PAREN},
		{"{", TOKEN_LEFT_BRACE},
		{"}", TOKEN_RIGHT_BRACE},
		{",", TOKEN_COMMA},
		{".", TOKEN_DOT},
		{"-", TOKEN_MINUS},
		{"+", TOKEN_PLUS},
		{";", TOKEN_SEMICOLON},
		{"/", TOKEN_SLASH},
		{"*", TOKEN_STAR},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			scanner := newScanner(test.input)
			result := scanner.scanToken()
			if result.tokenType != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result.tokenType)
			}
			if result.lexeme != test.input {
				t.Errorf("Expected lexeme %v, got %v", test.input, result.lexeme)
			}
		})
	}

	t.Run("All tokens", func(t *testing.T) {
		source := "( ) { } , . - + ; / * //hello \n ("
		scanner := newScanner(source)
		for _, expected := range []TokenType{
			TOKEN_LEFT_PAREN, TOKEN_RIGHT_PAREN,
			TOKEN_LEFT_BRACE, TOKEN_RIGHT_BRACE,
			TOKEN_COMMA, TOKEN_DOT,
			TOKEN_MINUS, TOKEN_PLUS,
			TOKEN_SEMICOLON, TOKEN_SLASH,
			TOKEN_STAR, TOKEN_LEFT_PAREN, TOKEN_EOF,
		} {
			result := scanner.scanToken()
			if result.tokenType != expected {
				t.Errorf("Expected %v, got %v", expected, result.tokenType)
			}
		}
	})

	t.Run("String token", func(t *testing.T) {
		source := `"hello" `
		scanner := newScanner(source)
		result := scanner.scanToken()
		if result.tokenType != TOKEN_STRING {
			t.Errorf("Expected TOKEN_STRING, got %v", result.tokenType)
		}
		if result.lexeme != `"hello"` {
			t.Errorf("Expected lexeme %v, got %v", `"hello"`, result.lexeme)
		}
	})

	t.Run("Number token", func(t *testing.T) {
		source := "12345 "
		scanner := newScanner(source)
		result := scanner.scanToken()
		if result.tokenType != TOKEN_NUMBER {
			t.Errorf("Expected TOKEN_NUMBER, got %v", result.tokenType)
		}
		if result.lexeme != "12345" {
			t.Errorf("Expected lexeme %v, got %v", "12345", result.lexeme)
		}
	})

	keywordsAndIdentifiers := []struct {
		input    string
		expected TokenType
	}{
		{
			"if ",
			TOKEN_IF,
		},
		{
			"else ",
			TOKEN_ELSE,
		},
		{
			" hello",
			TOKEN_IDENTIFIER,
		},
		{
			"superb",
			TOKEN_IDENTIFIER,
		},
		{
			"super",
			TOKEN_SUPER,
		},
		{
			"for",
			TOKEN_FOR,
		},
		{
			"fort",
			TOKEN_IDENTIFIER,
		},
	}

	for _, test := range keywordsAndIdentifiers {
		t.Run(test.input, func(t *testing.T) {
			scanner := newScanner(test.input)
			result := scanner.scanToken()
			if result.tokenType != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result.tokenType)
			}
			if result.lexeme != strings.TrimSpace(test.input) {
				t.Errorf("Expected lexeme %v, got %v", strings.TrimSpace(test.input), result.lexeme)
			}
		})
	}
}
