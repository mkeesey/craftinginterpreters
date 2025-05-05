package main

import (
	"testing"

	"github.com/mkeesey/craftinginterpreters/pkg/bytecode"
)

func TestInterpret(t *testing.T) {
	vm := bytecode.NewVM()
	defer vm.Free()

	source := `1 + 1`

	err := vm.Interpret(source)
	if err != nil {
		t.Fatalf("Interpret failed: %v", err)
	}
}
