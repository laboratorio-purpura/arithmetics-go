// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Add computes the sum of two integers.
//
// Add stores into sum the len(sum) least significant words of the result.
func Add(sum, x, y []uint) (carry uint) {
	// ensure x is not shorter than y
	if len(x) < len(y) {
		x, y = y, x
	}

	// truncate result to len(r) words
	xz := min(len(sum), len(x))
	yz := min(len(sum), len(y))
	// invariant: len(r) >= xz >= yz

	// add words, propagating carry
	for i := 0; i < yz; i++ {
		sum[i], carry = bits.Add(x[i], y[i], carry)
	}

	// propagate carry
	for i := yz; i < xz; i++ {
		sum[i], carry = bits.Add(x[i], 0, carry)
	}

	return
}
