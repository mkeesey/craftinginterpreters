package failure

import (
	"fmt"
	"strings"

	"github.com/mkeesey/craftinginterpreters/pkg/token"
)

func Error(line int, message string) error {
	return Report(line, "", message)
}

func TokenError(tok *token.Token, message string) error {
	if tok.Type == token.EOF {
		return Report(tok.Line, "at end", message)
	} else {
		return Report(tok.Line, fmt.Sprintf(" at '%s'", tok.Lexeme), message)
	}
}

func Report(line int, where string, message string) error {
	whereStr := strings.TrimSuffix(where, "\n")
	return fmt.Errorf("[%d] Error %s: %s", line, whereStr, message)
}

func Wrap(line int, message string, err error) error {
	return fmt.Errorf("[%d] Error %s: %w", line, message, err)
}
