package arithmetics

import "math/bits"

// Halve computes an integer half to a power.
//
// Halve stores into quotient the len(quotient) least significant words of the result.
// It permits aliasing quotient to x, in which case it becomes "halve accumulate".
//
// This implementation applies the "binary shift" method.
func Halve(quotient []uint, x []uint, y uint) (remainder uint) {
	const Bits = bits.UintSize

	qz := len(quotient)
	xz := len(x)

	// TODO: lift this restriction
	if y >= Bits {
		panic("y >= Bits")
	}

	// count of result words to compute
	z := min(qz, xz)

	// halve word by word,
	// from most to least significant,
	// propagating remainder
	for i := z; i > 0; i-- {
		// x[i] ÷ 2^y
		q := x[i-1] >> y
		r := x[i-1] << (Bits - y)
		// store quotient, propagate remainder
		quotient[i-1] = q + remainder
		remainder = r
	}

	return
}
