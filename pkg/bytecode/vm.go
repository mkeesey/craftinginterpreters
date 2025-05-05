package bytecode

import (
	"fmt"
	"os"
)

var Debug = false

var ErrInterpretError = fmt.Errorf("interpret error")
var ErrRuntimeError = fmt.Errorf("runtime error")
var InterpretRuntimeError = fmt.Errorf("interpret runtime error")

const StackMax = 256

type VM struct {
	chunk    *Chunk
	stack    [StackMax]Value
	ip       int
	stackIdx int
}

func NewVM() *VM {
	vm := &VM{
		chunk:    nil,
		stack:    [StackMax]Value{},
		ip:       0,
		stackIdx: 0,
	}
	vm.resetStack()
	return vm
}

func (vm *VM) resetStack() {
	vm.stackIdx = 0
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
				fmt.Printf("[ %s ]", vm.stack[i])
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
		case OP_NIL:
			vm.push(NilValue())
		case OP_TRUE:
			vm.push(BoolValue(true))
		case OP_FALSE:
			vm.push(BoolValue(false))
		case OP_EQUAL:
			b := vm.pop()
			a := vm.pop()
			vm.push(BoolValue(valuesEqual(a, b)))
		case OP_GREATER:
			vm.binaryOp(greater)
		case OP_LESS:
			vm.binaryOp(less)
		case OP_ADD:
			vm.binaryOp(add)
		case OP_SUBTRACT:
			vm.binaryOp(subtract)
		case OP_MULTIPLY:
			vm.binaryOp(multiply)
		case OP_DIVIDE:
			vm.binaryOp(divide)
		case OP_NOT:
			vm.push(isFalsy(vm.pop()))
		case OP_NEGATE:
			if !vm.peek(0).IsNumber() {
				vm.runtimeError("Operand must be a number.")
				return InterpretRuntimeError
			}
			vm.push(NumberValue(-(vm.pop().AsNumber())))
		case OP_RETURN:
			fmt.Printf("%s\n", vm.pop())
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

func (vm *VM) peek(distance int) Value {
	if vm.stackIdx == 0 {
		panic("Stack underflow")
	}
	return vm.stack[vm.stackIdx-1-distance]
}

func (vm *VM) binaryOp(op func(a, b float64) Value) error {
	if !vm.peek(0).IsNumber() || !vm.peek(1).IsNumber() {
		vm.runtimeError("Operands must be numbers.")
		return InterpretRuntimeError
	}
	b := vm.pop().AsNumber()
	a := vm.pop().AsNumber()
	vm.push(op(a, b))
	return nil
}

func (vm *VM) runtimeError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")

	line := vm.chunk.lines[vm.ip-1]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	vm.resetStack()
}

func isFalsy(value Value) Value {
	if value.IsNil() {
		return BoolValue(true)
	}
	if value.IsBool() {
		return BoolValue(!value.AsBool())
	}
	return BoolValue(false)
}
