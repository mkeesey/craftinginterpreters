package scanner

import (
	"strings"
	"testing"

	"github.com/mkeesey/craftinginterpreters/token"
	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	t.Run("single chars only", func(t *testing.T) {
		scanner := NewScanner(strings.NewReader("*+("))
		tokens, err := scanner.scanTokens()
		require.Nil(t, err)
		require.Len(t, tokens, 4)
		require.Equal(t, token.NewToken(token.STAR, "*", nil, 0), tokens[0])
		require.Equal(t, token.NewToken(token.PLUS, "+", nil, 0), tokens[1])
		require.Equal(t, token.NewToken(token.LEFT_PAREN, "(", nil, 0), tokens[2])
		require.Equal(t, token.NewToken(token.EOF, "", nil, 0), tokens[3])
	})
	t.Run("unknown char", func(t *testing.T) {
		scanner := NewScanner(strings.NewReader("*$-"))
		tokens, err := scanner.scanTokens()
		require.NotNil(t, err) // unknown $
		require.Len(t, tokens, 3)
		require.Equal(t, token.NewToken(token.STAR, "*", nil, 0), tokens[0])
		require.Equal(t, token.NewToken(token.MINUS, "-", nil, 0), tokens[1])
		require.Equal(t, token.NewToken(token.EOF, "", nil, 0), tokens[2])
	})

	t.Run("bang noequal", func(t *testing.T) {
		scanner := NewScanner(strings.NewReader("!*("))
		tokens, err := scanner.scanTokens()
		require.Nil(t, err)
		require.Len(t, tokens, 4)
		require.Equal(t, token.NewToken(token.BANG, "!", nil, 0), tokens[0])
		require.Equal(t, token.NewToken(token.STAR, "*", nil, 0), tokens[1])
		require.Equal(t, token.NewToken(token.LEFT_PAREN, "(", nil, 0), tokens[2])
		require.Equal(t, token.NewToken(token.EOF, "", nil, 0), tokens[3])
	})
	t.Run("bang equal", func(t *testing.T) {
		scanner := NewScanner(strings.NewReader("!=("))
		tokens, err := scanner.scanTokens()
		require.Nil(t, err)
		require.Len(t, tokens, 3)
		require.Equal(t, token.NewToken(token.BANG_EQUAL, "!=", nil, 0), tokens[0])
		require.Equal(t, token.NewToken(token.LEFT_PAREN, "(", nil, 0), tokens[1])
		require.Equal(t, token.NewToken(token.EOF, "", nil, 0), tokens[2])
	})

	t.Run("slash", func(t *testing.T) {
		type testcase struct {
			title  string
			input  string
			output []token.Token
		}

		testcases := []testcase{
			{
				title: "single slash",
				input: "/*",
				output: []token.Token{
					token.NewToken(token.SLASH, "/", nil, 0),
					token.NewToken(token.STAR, "*", nil, 0),
					token.NewToken(token.EOF, "", nil, 0),
				},
			},
			{
				title: "single line comment",
				input: " * //this is a comment",
				output: []token.Token{
					token.NewToken(token.STAR, "*", nil, 0),
					token.NewToken(token.EOF, "", nil, 0),
				},
			},
			{
				title: "single line comment without text",
				input: " * //",
				output: []token.Token{
					token.NewToken(token.STAR, "*", nil, 0),
					token.NewToken(token.EOF, "", nil, 0),
				},
			},
			{
				title: "multi line comment",
				input: "* //this is a comment\n (",
				output: []token.Token{
					token.NewToken(token.STAR, "*", nil, 0),
					token.NewToken(token.LEFT_PAREN, "(", nil, 1),
					token.NewToken(token.EOF, "", nil, 1),
				},
			},
		}
		for _, testcase := range testcases {
			t.Run(testcase.title, func(t *testing.T) {

				scanner := NewScanner(strings.NewReader(testcase.input))
				tokens, err := scanner.scanTokens()
				require.Nil(t, err)
				require.Len(t, tokens, len(testcase.output))
				for i, expected := range testcase.output {
					require.Equal(t, expected, tokens[i])
				}
			})
		}
	})

	t.Run("string literals", func(t *testing.T) {
		type testcase struct {
			title       string
			input       string
			expectError bool
			output      []token.Token
		}

		testcases := []testcase{
			{
				title:       "simple string",
				input:       `"hello"`,
				expectError: false,
				output: []token.Token{
					token.NewToken(token.STRING, "hello", "hello", 0),
					token.NewToken(token.EOF, "", nil, 0),
				},
			},
			{
				title: "more complicated string",
				input: `"hello
world"`,
				expectError: false,
				output: []token.Token{
					token.NewToken(token.STRING, "hello\nworld", "hello\nworld", 1),
					token.NewToken(token.EOF, "", nil, 1),
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
				scanner := NewScanner(strings.NewReader(testcase.input))
				tokens, err := scanner.scanTokens()
				if testcase.expectError {
					require.NotNil(t, err)
				} else {
					require.Nil(t, err)
					require.Len(t, tokens, len(testcase.output))
					for i, expected := range testcase.output {
						require.Equal(t, expected, tokens[i])
					}
				}
			})
		}
	})
}
