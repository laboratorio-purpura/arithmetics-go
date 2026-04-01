package arithmetics

import "math/bits"

// IsGreater tests if x is greater than y.
func IsGreater(x, y []uint) bool {
	xz := len(x)
	yz := len(y)

	// count of common words
	z := min(xz, yz)

	// x > y <=> y < x <=> x - y < 0
	// compare words from least to most significant,
	// computing the "borrow" of their difference.

	var borrow uint

	// common words
	for i := 0; i != z; i++ {
		_, borrow = bits.Sub(y[i], x[i], borrow)
	}

	// x excess words
	for i := z; i != xz; i++ {
		_, borrow = bits.Sub(0, x[i], borrow)
	}

	// y excess words
	for i := z; i != yz; i++ {
		_, borrow = bits.Sub(y[i], 0, borrow)
	}

	// borrow > 0 => y < x
	return borrow > 0
}

// NotGreater tests if x is *not* greater than y.
func NotGreater(x, y []uint) bool {
	return !IsGreater(x, y)
}
