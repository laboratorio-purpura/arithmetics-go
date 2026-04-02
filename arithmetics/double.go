// SPDX-FileCopyrightText: 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

// Double computes twice (to a power) of an integer.
//
// Double adds into product the len(product) least significant words of the result.
// It permits aliasing product to x, in which case it becomes "double accumulate".
//
// This implementation applies the "binary shift" method.
func Double(product []uint, x []uint, y uint) (excess uint) {
	const Bits = bits.UintSize

	pz := len(product)
	xz := len(x)

	// TODO: lift this restriction
	if y >= Bits {
		panic("y >= Bits")
	}

	// count of result words to compute
	z := min(pz, xz)

	// double word by word,
	// from least to most significant,
	// propagating excess
	for i := 0; i < z; i++ {
		// x[i] × 2^y
		p0 := x[i] << y
		p1 := x[i] >> (Bits - y)
		// store low word, propagate high word
		product[i] = p0 + excess
		excess = p1
	}

	return
}
