package bytecode

import "fmt"

var Debug = false

var ErrInterpretError = fmt.Errorf("interpret error")
var ErrRuntimeError = fmt.Errorf("runtime error")

const StackMax = 256

type VM struct {
	chunk    *Chunk
	stack    [StackMax]Value
	ip       int
	stackIdx int
}

func NewVM() *VM {
	return &VM{
		chunk:    nil,
		stack:    [StackMax]Value{},
		ip:       0,
		stackIdx: 0,
	}
}

func (vm *VM) Free() {
}

func (vm *VM) Interpret(source string) error {
	chunk := NewChunk()

	err := compile(source, chunk)
	if err != nil {
		return fmt.Errorf("Interpret: %w", err)
	}

	vm.chunk = chunk
	vm.ip = 0

	return vm.run()
}

func (vm *VM) run() error {
	for {
		if vm.ip >= len(vm.chunk.code) {
			return ErrInterpretError
		}

		if Debug {
			fmt.Printf("         ")
			for i := 0; i < vm.stackIdx; i++ {
				fmt.Printf("[ %g ]", vm.stack[i])
			}
			fmt.Printf("\n")
			disassembleInstruction(vm.chunk, vm.ip)
		}

		instruction := OpCode(readByte(vm.chunk.code, &vm.ip))

		switch instruction {
		case OP_CONSTANT:
			constantIndex := readByte(vm.chunk.code, &vm.ip)
			if constantIndex >= uint8(len(vm.chunk.constants)) {
				return ErrInterpretError
			}
			constant := vm.chunk.constants[constantIndex]
			vm.push(constant)
		case OP_ADD:
			vm.binaryOp(add)
		case OP_SUBTRACT:
			vm.binaryOp(subtract)
		case OP_MULTIPLY:
			vm.binaryOp(multiply)
		case OP_DIVIDE:
			vm.binaryOp(divide)
		case OP_NEGATE:
			vm.push(-vm.pop())
		case OP_RETURN:
			fmt.Printf("%g\n", vm.pop())
			return nil
		default:
			return ErrInterpretError
		}
	}
}

func readByte(code []byte, ip *int) byte {
	if *ip >= len(code) {
		return 0
	}
	b := code[*ip]
	*ip++
	return b
}

func (vm *VM) push(value Value) {
	if vm.stackIdx >= StackMax {
		panic("Stack overflow")
	}
	vm.stack[vm.stackIdx] = value
	vm.stackIdx++
}

func (vm *VM) pop() Value {
	if vm.stackIdx == 0 {
		panic("Stack underflow")
	}
	vm.stackIdx--
	return vm.stack[vm.stackIdx]
}

func (vm *VM) binaryOp(op func(a, b Value) Value) {
	b := vm.pop()
	a := vm.pop()
	vm.push(op(a, b))
}
