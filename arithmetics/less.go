// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// IsLess32 is true if and only if x is less than y.
func IsLess32(x, y []uint32) bool {
	xz := len(x)
	yz := len(y)
	z := min(xz, yz)

	var borrow uint32

	for i := 0; i != z; i++ {
		_, borrow = bits.Sub32(x[i], y[i], borrow)
	}

	for i := z; i != xz; i++ {
		_, borrow = bits.Sub32(x[i], 0, borrow)
	}

	for i := z; i != yz; i++ {
		_, borrow = bits.Sub32(0, y[i], borrow)
	}

	return borrow > 0
}

// NotLess32 is true if and only if x is not less than y.
func NotLess32(x, y []uint32) bool {
	return !IsLess32(x, y)
}

// IsMore32 is true if and only if x is more than y.
func IsMore32(x, y []uint32) bool {
	return IsLess32(y, x)
}

// NotMore32 is true if and only if x is not more than y.
func NotMore32(x, y []uint32) bool {
	return !IsLess32(y, x)
}
