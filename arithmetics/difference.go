// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Difference computes the difference between two integers.
//
// Difference stores into difference the len(difference) least significant words of the result.
// It permits aliasing difference to x, in which case it becomes "subtract accumulate".
//
// This implementation applies the "school" method described in Knuth, section 4.3.1.
func Difference(difference, x, y []uint) (borrow uint) {
	dz := len(difference)
	xz := len(x)
	yz := len(y)

	// count of result words to compute
	z := min(dz, xz, yz)

	// subtract x and y word by word,
	// from least to most significant,
	// propagating borrow
	for i := 0; i < z; i++ {
		difference[i], borrow = bits.Sub(x[i], y[i], borrow)
	}

	// either propagate borrow through x
	for i := z; i < min(dz, xz); i++ {
		difference[i], borrow = bits.Sub(x[i], 0, borrow)
	}

	// or propagate borrow through y
	for i := z; i < min(dz, yz); i++ {
		difference[i], borrow = bits.Sub(0, y[i], borrow)
	}

	return
}
