// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Sum of nonnegative integers `x` and `y`.
//
// Stores into `r` the `size(r)` least significant words of the result.
// Permits aliasing `r` to `x`, in which case it "accumulates" the result.
// Returns the "carry" of the top word of the result.
//
// This implementation applies the "school" method described in Knuth, section 4.3.1.
func Sum(r, x, y []uint) (c uint) {
	rz := len(r)
	xz := len(x)
	yz := len(y)

	// count of result words
	z := min(rz, xz, yz)

	// add x and y, word by word, propagating carry
	for i := 0; i < z; i++ {
		r[i], c = bits.Add(x[i], y[i], c)
	}

	// either propagate carry through x
	for i := z; i < min(rz, xz); i++ {
		r[i], c = bits.Add(x[i], 0, c)
	}

	// or propagate carry through y
	for i := z; i < min(rz, yz); i++ {
		r[i], c = bits.Add(0, y[i], c)
	}

	return
}
