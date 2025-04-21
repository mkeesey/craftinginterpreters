package bytecode

import (
	"errors"
	"fmt"
)

func DisassembleChunk(chunk *Chunk, name string) {
	fmt.Println("== ", name, " ==")
	var err error
	for offset := 0; offset < len(chunk.code); {
		offset, err = disassembleInstruction(chunk, offset)
		if err != nil {
			fmt.Printf("Error disassembling instruction at offset %d: %v\n", offset, err)
			break
		}
	}
}

func disassembleInstruction(chunk *Chunk, offset int) (int, error) {
	if offset >= len(chunk.code) {
		return offset, fmt.Errorf("offset %d out of bounds", offset)
	}
	fmt.Printf("%04d ", offset)
	if offset > 0 && chunk.lines[offset] == chunk.lines[offset-1] {
		fmt.Print("   | ")
	} else {
		fmt.Printf("%4d ", chunk.lines[offset])
	}

	instruction := OpCode(chunk.code[offset])
	switch instruction {
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", chunk, offset)
	case OP_ADD:
		return simpleInstruction("OP_ADD", offset), nil
	case OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset), nil
	case OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset), nil
	case OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset), nil
	case OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset), nil
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset), nil
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1, errors.New("unknown opcode")
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func constantInstruction(name string, chunk *Chunk, offset int) (int, error) {
	constantIndex := chunk.code[offset+1]
	if constantIndex >= uint8(len(chunk.constants)) {
		return offset, fmt.Errorf("constant index %d out of bounds", constantIndex)
	}
	constant := chunk.constants[constantIndex]
	fmt.Printf("%s %g\n", name, constant)
	return offset + 2, nil
}
