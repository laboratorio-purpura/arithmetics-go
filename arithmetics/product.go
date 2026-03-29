// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// Product32 of "long" integers.
//
// Computes r + x × y by the "school" method.
//
// Stores len(r) words of the result into r
func Product32(r, x, y []uint32) {
	xz := len(x)
	yz := len(y)

	for j := 0; j < yz; j++ {
		var excess uint32
		for i := 0; i < xz; i++ {
			var carry uint32
			// x[i] × y[j] + excess
			hi, lo := bits.Mul32(x[i], y[j])
			lo, carry = bits.Add32(lo, excess, 0)
			hi, _ = bits.Add32(hi, 0, carry)
			// store into r[i+j]
			r[i+j], carry = bits.Add32(r[i+j], lo, 0)
			// propagate excess
			excess, _ = bits.Add32(hi, 0, carry)
		}
		// store excess
		r[xz+j], _ = bits.Add32(r[xz+j], excess, 0)
	}
}
