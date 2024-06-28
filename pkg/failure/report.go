package failure

import (
	"fmt"
	"strings"
)

func Error(line int, message string) error {
	return Report(line, "", message)
}

func Report(line int, where string, message string) error {
	whereStr := strings.TrimSuffix(where, "\n")
	return fmt.Errorf("[%d] Error %s: %s", line, whereStr, message)
}

func Wrap(line int, message string, err error) error {
	return fmt.Errorf("[%d] Error %s: %w", line, message, err)
}
