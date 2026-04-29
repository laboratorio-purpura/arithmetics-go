// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Sum computes the sum of two integers.
//
// Sum stores into sum the len(sum) least significant words of the result.
// It permits aliasing sum to x, in which case it becomes "add accumulate".
//
// This implementation applies the "school" method described in Knuth, section 4.3.1.
func Sum(sum, x, y []uint) (carry uint) {
	sz := len(sum)
	xz := len(x)
	yz := len(y)

	// count of result words to compute
	z := min(sz, xz, yz)

	// add x and y word by word,
	// from least to most significant,
	// propagating carry
	for i := 0; i < z; i++ {
		sum[i], carry = bits.Add(x[i], y[i], carry)
	}

	// either propagate carry through x
	for i := z; i < min(sz, xz); i++ {
		sum[i], carry = bits.Add(x[i], 0, carry)
	}

	// or propagate carry through y
	for i := z; i < min(sz, yz); i++ {
		sum[i], carry = bits.Add(0, y[i], carry)
	}

	return
}
