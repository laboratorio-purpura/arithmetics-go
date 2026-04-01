// SPDX-FileCopyrightText: 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Double computes twice an integer, y times.
//
// Double computes x × 2^y.
//
// Double adds product into the len(product) least significant words of the result.
// If product starts nonzero, this operation becomes "double and add".
//
// In a binary machine, this is equivalent to a left shift by y.
func Double(product []uint, x []uint, y uint) (excess uint) {
	const Bits = bits.UintSize

	pz := len(product)
	xz := len(x)

	// TODO: lift this restriction
	if y >= Bits {
		panic("y >= Bits")
	}

	// compute count of resulting words
	z := min(pz, xz)

	// compute most significant resulting word
	if z > 0 {
		excess = x[z-1] >> (Bits - y)
		product[z-1] += x[z-1] << y
	}

	// compute remaining resulting words
	for i := z; i > 1; i-- {
		product[i-1] += x[i-2] >> (Bits - y)
		product[i-2] += x[i-2] << y
	}

	return
}
