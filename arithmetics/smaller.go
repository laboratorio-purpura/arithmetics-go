// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// IsSmaller is true if and only if x is less than y.
func IsSmaller(x, y []uint) bool {
	xz := len(x)
	yz := len(y)
	z := min(xz, yz)

	var borrow uint

	for i := 0; i != z; i++ {
		_, borrow = bits.Sub(x[i], y[i], borrow)
	}

	for i := z; i != xz; i++ {
		_, borrow = bits.Sub(x[i], 0, borrow)
	}

	for i := z; i != yz; i++ {
		_, borrow = bits.Sub(0, y[i], borrow)
	}

	return borrow > 0
}

// NotSmaller is true if and only if x is not less than y.
func NotSmaller(x, y []uint) bool {
	return !IsSmaller(x, y)
}
