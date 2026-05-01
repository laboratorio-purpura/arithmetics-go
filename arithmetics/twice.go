// SPDX-FileCopyrightText: 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Twice of nonnegative integer `x`, `y` times.
//
// Stores into `r` the `size(r)` least significant words of the result.
// Permits aliasing `r` to `x`, in which case it "accumulates" the result.
// Returns the "excess" of the top word of the result.
//
// This implementation applies the "shift" method.
func Twice(r []uint, x []uint, y uint) (e uint) {
	const Bits = bits.UintSize

	rz := len(r)
	xz := len(x)

	// TODO: document this restriction
	y = min(y, Bits-1)

	// count of result words
	z := min(rz, xz)

	// double x, y times, word by word, propagating excess
	for i := 0; i < z; i++ {
		// x[i] × 2^y
		p0 := x[i] << y
		p1 := x[i] >> (Bits - y)
		// store low word with excess, forward high word
		r[i] = p0 | e
		e = p1
	}

	return
}
