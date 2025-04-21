package main

import (
	"github.com/mkeesey/craftinginterpreters/pkg/bytecode"
)

func main() {
	vm := bytecode.NewVM()
	defer vm.Free()

	bytecode.Debug = true

	// Create a new chunk
	chunk := bytecode.NewChunk()

	constant := chunk.WriteConstant(1.2)
	chunk.Write(byte(bytecode.OP_CONSTANT), 123)
	chunk.Write(byte(constant), 123)

	constant = chunk.WriteConstant(3.4)
	chunk.Write(byte(bytecode.OP_CONSTANT), 123)
	chunk.Write(byte(constant), 123)

	chunk.Write(byte(bytecode.OP_ADD), 123)

	constant = chunk.WriteConstant(5.6)
	chunk.Write(byte(bytecode.OP_CONSTANT), 123)
	chunk.Write(byte(constant), 123)

	chunk.Write(byte(bytecode.OP_DIVIDE), 123)
	chunk.Write(byte(bytecode.OP_NEGATE), 123)

	// Write an OP_RETURN instruction to the chunk
	chunk.Write(byte(bytecode.OP_RETURN), 123)

	err := vm.Interpret(chunk)
	if err != nil {
		//TODO
		panic(err)
	}

	// Free the chunk's resources
	chunk.Free()
}
