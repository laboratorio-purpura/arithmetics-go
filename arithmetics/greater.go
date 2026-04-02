package arithmetics

// IsGreater tests if an integer is greater than another.
func IsGreater(x, y []uint) bool {
	xz := len(x)
	yz := len(y)

	z := min(xz, yz)

	for i := xz; i > z; i-- {
		if x[i-1] != 0 {
			return true
		}
	}

	for i := yz; i > z; i-- {
		if y[i-1] != 0 {
			return false
		}
	}

	for i := z; i > 0; i-- {
		if x[i-1] < y[i-1] {
			return false
		}
		if x[i-1] > y[i-1] {
			return true
		}
	}

	return false
}

// NotGreater tests if x is not greater than y.
func NotGreater(x, y []uint) bool {
	return !IsGreater(x, y)
}
