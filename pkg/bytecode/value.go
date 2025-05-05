package bytecode

import "fmt"

type ValueType int

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
)

func BoolValue(b bool) Value {
	return Value{Type: VAL_BOOL, Value: b}
}

func NilValue() Value {
	return Value{Type: VAL_NIL, Value: nil}
}

func NumberValue(n float64) Value {
	return Value{Type: VAL_NUMBER, Value: n}
}

type Value struct {
	Type  ValueType
	Value any
}

func (v Value) String() string {
	return fmt.Sprintf("%v", v.Value)
}

func (v Value) IsBool() bool {
	return v.Type == VAL_BOOL
}

func (v Value) IsNil() bool {
	return v.Type == VAL_NIL
}

func (v Value) IsNumber() bool {
	return v.Type == VAL_NUMBER
}

func (v Value) AsBool() bool {
	return v.Value.(bool)
}

func (v Value) AsNumber() float64 {
	return v.Value.(float64)
}

type ValueArray []Value

func valuesEqual(a, b Value) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case VAL_BOOL:
		return a.AsBool() == b.AsBool()
	case VAL_NIL:
		return true
	case VAL_NUMBER:
		return a.AsNumber() == b.AsNumber()
	default:
		return false
	}
}
