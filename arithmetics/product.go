// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// ProductUni of nonnegative integers `x` and `y`.
//
// Stores into `r` the `size(r)` least significant words of the result.
// Permits aliasing `r` to `x`, in which case it "accumulates" the result.
// Returns the "excess" of the top word of the result.
//
// This implementation applies the "school" method described in Knuth, section 4.3.1.
func ProductUni(r []uint, x []uint, y uint) (e uint) {
	rz := len(r)
	xz := len(x)

	// count of result words
	z := min(xz, rz)

	// multiply x and y, word by word, propagating excess
	for i := 0; i < z; i++ {
		// x[i] × y
		p1, p0 := bits.Mul(x[i], y)
		// add lower excess
		var carry uint
		p0, carry = bits.Add(p0, e, 0)
		p1, _ = bits.Add(p1, 0, carry)
		// store low word, forward high word
		r[i] = p0
		e = p1
	}

	return
}

// Product of nonnegative integers `x` and `y`.
//
// Stores into `r` the `size(r)` least significant words of the result.
//
// This implementation applies the "school" method described in Knuth, section 4.3.1.
func Product(r, x, y []uint) {
	rz := len(r)
	xz := len(x)
	yz := len(y)

	clear(r)

	// multiply x and y, word by word, propagating excess
	for j := 0; j < min(rz, yz); j++ {
		var excess uint
		for i := 0; i < min(rz, xz) && i+j < rz; i++ {
			var carry uint
			// x[i] × y[j]
			p1, p0 := bits.Mul(x[i], y[j])
			// add lower excess
			p0, carry = bits.Add(p0, excess, 0)
			p1, _ = bits.Add(p1, 0, carry)
			// store low word, forward high word
			r[i+j], carry = bits.Add(r[i+j], p0, 0)
			excess, _ = bits.Add(0, p1, carry)
		}
		// store excess
		if xz+j < rz {
			r[xz+j], _ = bits.Add(r[xz+j], excess, 0)
		}
	}
}
