package failure

import (
	"fmt"
	"os"
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
		return Report(tok.Line, fmt.Sprintf("at '%s'", tok.Lexeme), message)
	}
}

func Report(line int, where string, message string) error {
	whereStr := strings.TrimSuffix(where, "\n")
	return fmt.Errorf("[line: %d] Error %s: %s", line, whereStr, message)
}

func Wrap(line int, message string, err error) error {
	return fmt.Errorf("[line: %d] Error %s: %w", line, message, err)
}

type Reporter struct {
	hasFailed bool
}

func (r *Reporter) Error(line int, message string) {
	r.Report(line, "", message)
}

func (r *Reporter) TokenError(tok *token.Token, message string) {
	if tok.Type == token.EOF {
		r.Report(tok.Line, "at end", message)
	} else {
		r.Report(tok.Line, fmt.Sprintf("at '%s'", tok.Lexeme), message)
	}
}

func (r *Reporter) Report(line int, where string, message string) {
	whereStr := strings.TrimSuffix(where, "\n")
	fmt.Fprintf(os.Stderr, "[line: %d] Error %s: %s\n", line, whereStr, message)
	r.hasFailed = true
}

func (r *Reporter) ReportErr(line int, message string, err error) {
	fmt.Fprintf(os.Stderr, "[line: %d] Error %s: %s\n", line, message, err)
	r.hasFailed = true
}

func (r *Reporter) Panic(line int, err error) {
	r.hasFailed = true
	panic(fmt.Sprintf("line %d: %s", line, err))
}

func (r *Reporter) HasFailed() bool {
	return r.hasFailed
}

func (r *Reporter) Reset() {
	r.hasFailed = false
}
