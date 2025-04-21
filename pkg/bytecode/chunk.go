package bytecode

type OpCode byte

const (
	OP_CONSTANT OpCode = iota
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NEGATE
	OP_RETURN
)

type Chunk struct {
	code      []byte
	lines     []int
	constants ValueArray
}

func NewChunk() *Chunk {
	return &Chunk{
		code:      make([]byte, 0),
		lines:     make([]int, 0),
		constants: make(ValueArray, 0),
	}
}

func (c *Chunk) Write(b byte, line int) {
	c.code = append(c.code, b)
	c.lines = append(c.lines, line)
}

func (c *Chunk) WriteConstant(value Value) uint8 {
	c.constants = append(c.constants, value)
	return uint8(len(c.constants) - 1)
}

func (c *Chunk) Free() {
	c.code = c.code[:0]
	c.lines = c.lines[:0]
	c.constants = c.constants[:0]
}
