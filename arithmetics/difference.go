// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Difference32 with borrow of "long" integers.
//
// Computes x − y by the "school" method.
//
// Stores len(r) words of the result into r;
// returns the borrow.
func Difference32(r, x, y []uint32) uint32 {
	xz := min(len(r), len(x))
	yz := min(len(r), len(y))

	var borrow uint32

	// subtract words, propagating borrow
	z := min(xz, yz)
	for i := 0; i < z; i++ {
		var difference uint32
		difference, borrow = bits.Sub32(x[i], y[i], borrow)
		r[i] = difference
	}

	// propagate borrow
	if xz > yz {
		// y is shorter
		for i := yz; i < xz; i++ {
			var difference uint32
			difference, borrow = bits.Sub32(x[i], 0, borrow)
			r[i] = difference
		}
	} else {
		// y is not shorter
		for i := xz; i < yz; i++ {
			var difference uint32
			difference, borrow = bits.Sub32(0, y[i], borrow)
			r[i] = difference
		}
	}

	// invariant: borrow == 0 || borrow == 1
	return borrow
}
