package arithmetics

import "math/bits"

// Half of nonnegative integer `x`, `y` times.
//
// Stores into `q` the `size(q)` least significant words of the quotient.
// Returns the remainder.
//
// This implementation applies the "shift" method.
func Half(q []uint, x []uint, y uint) (r uint) {
	const Bits = bits.UintSize

	qz := len(q)
	xz := len(x)

	// TODO: document this restriction
	y = min(y, Bits-1)

	// count of result words
	z := min(qz, xz)

	// halve x, y times, word by word, propagating remainder
	if xz > z {
		// x[z] ÷ 2^y
		// forward remainder, preshifted
		r = x[z] << (Bits - y)
	}
	for i := z; i > 0; i-- {
		// x[i] ÷ 2^y
		q_ := x[i-1] >> y
		r_ := x[i-1] << (Bits - y)
		// store quotient with lower remainder
		q[i-1] = q_ | r
		// forward current remainder, preshifted
		r = r_
	}

	return
}
