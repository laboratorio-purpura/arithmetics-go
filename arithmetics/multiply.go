// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// Multiply computes the product of two integers.
//
// Multiply adds into product the len(product) least significant words of the result.
// If product starts nonzero, this operation becomes "multiply and add".
//
// Computes by the "school" method.
func Multiply(product, x, y []uint) {
	xz := len(x)
	yz := len(y)

	for j := 0; j < yz; j++ {
		var excess uint
		for i := 0; i < xz; i++ {
			var carry uint
			// x[i] × y[j] + excess
			hi, lo := bits.Mul(x[i], y[j])
			lo, carry = bits.Add(lo, excess, 0)
			hi, _ = bits.Add(hi, 0, carry)
			// store into r[i+j]
			product[i+j], carry = bits.Add(product[i+j], lo, 0)
			// propagate excess
			excess, _ = bits.Add(hi, 0, carry)
		}
		// store excess
		product[xz+j], _ = bits.Add(product[xz+j], excess, 0)
	}
}
