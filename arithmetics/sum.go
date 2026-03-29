// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Sum32 with carry of "long" integers.
//
// Computes x + y by the "school" method.
//
// Stores len(r) words of the result into r;
// returns the carry.
func Sum32(r, x, y []uint32) uint32 {
	// ensure x is not shorter than y
	if len(x) < len(y) {
		x, y = y, x
	}

	// truncate result to len(r) words
	xz := min(len(r), len(x))
	yz := min(len(r), len(y))
	// invariant: len(r) >= xz >= yz

	var carry uint32

	// add words, propagating carry
	for i := 0; i < yz; i++ {
		var sum uint32
		sum, carry = bits.Add32(x[i], y[i], carry)
		r[i] = sum
	}

	// propagate carry
	for i := yz; i < xz; i++ {
		var sum uint32
		sum, carry = bits.Add32(x[i], 0, carry)
		r[i] = sum
	}

	// invariant: carry == 0 || carry == 1
	return carry
}
