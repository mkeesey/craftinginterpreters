package bytecode

func greater(a, b float64) Value {
	return BoolValue(a > b)
}

func less(a, b float64) Value {
	return BoolValue(a < b)
}

func add(a, b float64) Value {
	return NumberValue(a + b)
}

func subtract(a, b float64) Value {
	return NumberValue(a - b)
}

func multiply(a, b float64) Value {
	return NumberValue(a * b)
}

func divide(a, b float64) Value {
	return NumberValue(a / b)
}
