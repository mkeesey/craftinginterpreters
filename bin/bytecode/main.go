package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mkeesey/craftinginterpreters/pkg/bytecode"
)

func main() {
	vm := bytecode.NewVM()
	defer vm.Free()

	if len(os.Args) == 1 {
		repl(vm)
	} else if len(os.Args) == 2 {
		runFile(vm, os.Args[1])
	} else {
		fmt.Println("Usage: clox [path]")
		os.Exit(64)
	}
}

func repl(vm *bytecode.VM) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("failed to read line: %w", err)
		}
		vm.Interpret(line)
	}
}

func runFile(vm *bytecode.VM, path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(74)
	}
	defer file.Close()

	sourceBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(74)
	}

	if err := vm.Interpret(string(sourceBytes)); err != nil {
		fmt.Fprintf(os.Stderr, "Error interpreting file: %v\n", err)
		os.Exit(70)
	}
}
