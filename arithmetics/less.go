// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// IsLess is true if and only if x is less than y.
func IsLess(x, y []uint) bool {
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

// NotLess is true if and only if x is not less than y.
func NotLess(x, y []uint) bool {
	return !IsLess(x, y)
}

// IsMore is true if and only if x is more than y.
func IsMore(x, y []uint) bool {
	return IsLess(y, x)
}

// NotMore is true if and only if x is not more than y.
func NotMore(x, y []uint) bool {
	return !IsLess(y, x)
}
