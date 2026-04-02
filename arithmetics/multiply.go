// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// Multiply computes the product of two integers.
//
// Multiply adds into product the len(product) least significant words of the result.
// It forbids aliasing product to x.
// It permits product to start nonzero, in which case it becomes "multiply and add".
//
// This implementation applies the "school" method described in Knuth, section 4.3.1.
func Multiply(product, x, y []uint) {
	xz := len(x)
	yz := len(y)

	// TODO: lift this restriction
	if len(product) < xz+yz {
		panic("product is too small")
	}

	// multiply x and y word by word,
	// from least to most significant
	for j := 0; j < yz; j++ {
		var excess uint
		for i := 0; i < xz; i++ {
			var carry uint
			// x[i] × y[j] + excess
			p1, p0 := bits.Mul(x[i], y[j])
			p0, carry = bits.Add(p0, excess, 0)
			p1, _ = bits.Add(p1, 0, carry)
			// store low word, propagate high word
			product[i+j], carry = bits.Add(product[i+j], p0, 0)
			excess, _ = bits.Add(p1, 0, carry)
		}
		// store excess
		product[xz+j], _ = bits.Add(product[xz+j], excess, 0)
	}
}
