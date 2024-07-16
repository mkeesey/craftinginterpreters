package scanner

import (
	"strings"
	"testing"

	"github.com/mkeesey/craftinginterpreters/pkg/failure"
	"github.com/mkeesey/craftinginterpreters/pkg/token"
	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	t.Run("single chars only", func(t *testing.T) {
		reporter := &failure.Reporter{}
		scanner := NewScanner(strings.NewReader("*+("), reporter)
		tokens := scanner.ScanTokens()
		require.Len(t, tokens, 4)
		require.Equal(t, token.NewToken(token.STAR, "*", nil, 1), tokens[0])
		require.Equal(t, token.NewToken(token.PLUS, "+", nil, 1), tokens[1])
		require.Equal(t, token.NewToken(token.LEFT_PAREN, "(", nil, 1), tokens[2])
		require.Equal(t, token.NewToken(token.EOF, "", nil, 1), tokens[3])
	})
	t.Run("unknown char", func(t *testing.T) {
		reporter := &failure.Reporter{}
		scanner := NewScanner(strings.NewReader("*$-"), reporter)
		tokens := scanner.ScanTokens()
		require.True(t, reporter.HasFailed())
		require.Len(t, tokens, 3)
		require.Equal(t, token.NewToken(token.STAR, "*", nil, 1), tokens[0])
		require.Equal(t, token.NewToken(token.MINUS, "-", nil, 1), tokens[1])
		require.Equal(t, token.NewToken(token.EOF, "", nil, 1), tokens[2])
	})

	t.Run("bang noequal", func(t *testing.T) {
		reporter := &failure.Reporter{}
		scanner := NewScanner(strings.NewReader("!*("), reporter)
		tokens := scanner.ScanTokens()
		require.False(t, reporter.HasFailed())
		require.Len(t, tokens, 4)
		require.Equal(t, token.NewToken(token.BANG, "!", nil, 1), tokens[0])
		require.Equal(t, token.NewToken(token.STAR, "*", nil, 1), tokens[1])
		require.Equal(t, token.NewToken(token.LEFT_PAREN, "(", nil, 1), tokens[2])
		require.Equal(t, token.NewToken(token.EOF, "", nil, 1), tokens[3])
	})
	t.Run("bang equal", func(t *testing.T) {
		reporter := &failure.Reporter{}
		scanner := NewScanner(strings.NewReader("!=("), reporter)
		tokens := scanner.ScanTokens()
		require.False(t, reporter.HasFailed())
		require.Len(t, tokens, 3)
		require.Equal(t, token.NewToken(token.BANG_EQUAL, "!=", nil, 1), tokens[0])
		require.Equal(t, token.NewToken(token.LEFT_PAREN, "(", nil, 1), tokens[1])
		require.Equal(t, token.NewToken(token.EOF, "", nil, 1), tokens[2])
	})

	t.Run("slash", func(t *testing.T) {
		type testcase struct {
			title  string
			input  string
			output []*token.Token
		}

		testcases := []testcase{
			{
				title: "single slash",
				input: "/*",
				output: []*token.Token{
					token.NewToken(token.SLASH, "/", nil, 1),
					token.NewToken(token.STAR, "*", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title: "single line comment",
				input: " * //this is a comment",
				output: []*token.Token{
					token.NewToken(token.STAR, "*", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title: "single line comment without text",
				input: " * //",
				output: []*token.Token{
					token.NewToken(token.STAR, "*", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title: "multi line comment",
				input: "* //this is a comment\n (",
				output: []*token.Token{
					token.NewToken(token.STAR, "*", nil, 1),
					token.NewToken(token.LEFT_PAREN, "(", nil, 2),
					token.NewToken(token.EOF, "", nil, 2),
				},
			},
		}
		for _, testcase := range testcases {
			t.Run(testcase.title, func(t *testing.T) {

				reporter := &failure.Reporter{}
				scanner := NewScanner(strings.NewReader(testcase.input), reporter)
				tokens := scanner.ScanTokens()
				require.False(t, reporter.HasFailed())
				validateTokens(t, tokens, testcase.output)
			})
		}
	})

	t.Run("string literals", func(t *testing.T) {
		type testcase struct {
			title       string
			input       string
			expectError bool
			output      []*token.Token
		}

		testcases := []testcase{
			{
				title:       "simple string",
				input:       `"hello"`,
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.STRING, "hello", "hello", 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title: "more complicated string",
				input: `"hello
world"`,
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.STRING, "hello\nworld", "hello\nworld", 2),
					token.NewToken(token.EOF, "", nil, 2),
				},
			},
			{
				title:       "unterminated string",
				input:       `"hello`,
				expectError: true,
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.title, func(t *testing.T) {
				reporter := &failure.Reporter{}
				scanner := NewScanner(strings.NewReader(testcase.input), reporter)
				tokens := scanner.ScanTokens()
				validateResp(t, testcase.expectError, reporter, tokens, testcase.output)
			})
		}
	})

	t.Run("numbers", func(t *testing.T) {
		type testcase struct {
			title       string
			input       string
			expectError bool
			output      []*token.Token
		}

		testcases := []testcase{
			{
				title:       "single digit",
				input:       "1",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.NUMBER, "1", 1.0, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title:       "simple number",
				input:       "123",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.NUMBER, "123", 123.0, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title:       "decimal",
				input:       "123.456",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.NUMBER, "123.456", 123.456, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title:       "trailing dot",
				input:       "123.",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.NUMBER, "123", 123.0, 1),
					token.NewToken(token.DOT, ".", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title:       "dot without number",
				input:       `123."hello"`,
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.NUMBER, "123", 123.0, 1),
					token.NewToken(token.DOT, ".", nil, 1),
					token.NewToken(token.STRING, "hello", "hello", 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.title, func(t *testing.T) {
				reporter := &failure.Reporter{}
				scanner := NewScanner(strings.NewReader(testcase.input), reporter)
				tokens := scanner.ScanTokens()
				validateResp(t, testcase.expectError, reporter, tokens, testcase.output)
			})
		}
	})

	t.Run("identifiers", func(t *testing.T) {

		type testcase struct {
			title       string
			input       string
			expectError bool
			output      []*token.Token
		}

		testcases := []testcase{
			{
				title:       "simple identifier",
				input:       "hello",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.IDENTIFIER, "hello", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title:       "multiple identifier",
				input:       "hello world",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.IDENTIFIER, "hello", nil, 1),
					token.NewToken(token.IDENTIFIER, "world", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title:       "identifier followed by dot",
				input:       "hello.world",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.IDENTIFIER, "hello", nil, 1),
					token.NewToken(token.DOT, ".", nil, 1),
					token.NewToken(token.IDENTIFIER, "world", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title:       "identifier with underscore",
				input:       "hello_world",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.IDENTIFIER, "hello_world", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title:       "identifier complex",
				input:       "_he9llo_world",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.IDENTIFIER, "_he9llo_world", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
			{
				title:       "identifier and keywords",
				input:       "_he9llo_world or waffles",
				expectError: false,
				output: []*token.Token{
					token.NewToken(token.IDENTIFIER, "_he9llo_world", nil, 1),
					token.NewToken(token.OR, "or", nil, 1),
					token.NewToken(token.IDENTIFIER, "waffles", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
		}

		for _, testcase := range testcases {
			t.Run(testcase.title, func(t *testing.T) {
				reporter := &failure.Reporter{}
				scanner := NewScanner(strings.NewReader(testcase.input), reporter)
				tokens := scanner.ScanTokens()
				validateResp(t, testcase.expectError, reporter, tokens, testcase.output)
			})
		}
	})
}

func validateResp(t *testing.T, expectFailedReporter bool, reporter *failure.Reporter, tokens []*token.Token, expected []*token.Token) {
	t.Helper()
	if expectFailedReporter {
		require.True(t, reporter.HasFailed())
	} else {
		require.False(t, reporter.HasFailed())
		validateTokens(t, tokens, expected)
	}
}

func validateTokens(t *testing.T, tokens []*token.Token, expected []*token.Token) {
	t.Helper()
	require.Len(t, tokens, len(expected))
	for i, expected := range expected {
		require.Equal(t, expected, tokens[i])
	}
}
