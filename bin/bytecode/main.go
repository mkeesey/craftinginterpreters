package main

import (
	"github.com/mkeesey/craftinginterpreters/pkg/bytecode"
)

func main() {
	// Create a new chunk
	chunk := bytecode.NewChunk()

	constant := chunk.WriteConstant(1.2)
	chunk.Write(byte(bytecode.OP_CONSTANT), 123)
	chunk.Write(byte(constant), 123)

	// Write an OP_RETURN instruction to the chunk
	chunk.Write(byte(bytecode.OP_RETURN), 123)

	// Disassemble the chunk
	bytecode.DisassembleChunk(chunk, "test chunk")

	// Free the chunk's resources
	chunk.Free()
}
