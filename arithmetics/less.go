// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "golang.org/x/exp/constraints"

// IsLess is true if and only if x is less than y.
func IsLess[Word constraints.Unsigned](x, y []Word) bool {
	// operation is not commutative
	xz := len(x)
	yz := len(y)
	if xz == 0 && yz == 0 {
		return false
	}

	// test x excess digits: nonzero => x > y
	for i := xz - 1; i >= yz; i-- {
		if x[i] != 0 {
			return false
		}
	}

	// test y excess digits: nonzero => x < y
	for i := yz - 1; i >= xz; i-- {
		if y[i] != 0 {
			return true
		}
	}

	z := min(xz, yz)
	if z == 0 {
		return false
	}

	// test common digits from highest to lowest
	for i := z - 1; i >= 0; i-- {
		if x[i] >= y[i] {
			return false
		}
	}

	return true
}

// NotLess is true if and only if x is not less than y.
func NotLess[Word constraints.Unsigned](x, y []Word) bool {
	return !IsLess(x, y)
}

// IsMore is true if and only if x is more than y.
func IsMore[Word constraints.Unsigned](x, y []Word) bool {
	return IsLess(y, x)
}

// NotMore is true if and only if x is not more than y.
func NotMore[Word constraints.Unsigned](x, y []Word) bool {
	return !IsLess(y, x)
}
