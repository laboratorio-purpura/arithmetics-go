// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Difference of nonnegative integers `x` and `y`.
//
// Stores into `r` the `size(r)` least significant words of the result.
// Permits aliasing `r` to `x`, in which case it "accumulates" the result.
// Returns the "borrow" of the top word of the result.
//
// This implementation applies the "school" method described in Knuth, section 4.3.1.
func Difference(r, x, y []uint) (b uint) {
	rz := len(r)
	xz := len(x)
	yz := len(y)

	// count of result words
	z := min(rz, xz, yz)

	// subtract x and y, word by word, propagating borrow
	for i := 0; i < z; i++ {
		r[i], b = bits.Sub(x[i], y[i], b)
	}

	// either propagate borrow through x
	for i := z; i < min(rz, xz); i++ {
		r[i], b = bits.Sub(x[i], 0, b)
	}

	// or propagate borrow through y
	for i := z; i < min(rz, yz); i++ {
		r[i], b = bits.Sub(0, y[i], b)
	}

	return
}
