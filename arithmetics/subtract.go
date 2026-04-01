// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Subtract computes the difference between two integers.
//
// Subtract stores into difference the len(difference) least significant words of the result.
func Subtract(difference, x, y []uint) (borrow uint) {
	xz := min(len(difference), len(x))
	yz := min(len(difference), len(y))

	// subtract words, propagating borrow
	z := min(xz, yz)
	for i := 0; i < z; i++ {
		difference[i], borrow = bits.Sub(x[i], y[i], borrow)
	}

	// propagate borrow
	if xz > yz {
		// y is shorter
		for i := yz; i < xz; i++ {
			difference[i], borrow = bits.Sub(x[i], 0, borrow)
		}
	} else {
		// y is not shorter
		for i := xz; i < yz; i++ {
			difference[i], borrow = bits.Sub(0, y[i], borrow)
		}
	}

	// invariant: borrow == 0 || borrow == 1
	return borrow
}
