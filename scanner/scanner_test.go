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
}
