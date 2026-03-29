// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "golang.org/x/exp/constraints"

// AreEqual is true if and only if x and y are equal.
func AreEqual[Word constraints.Unsigned](x, y []Word) bool {
	// ensure x is not shorter than y
	if len(x) < len(y) {
		x, y = y, x
	}
	xz := len(x)
	yz := len(y)
	// invariant: xz >= yz

	// test common digits in whatever order
	for i := 0; i < yz; i++ {
		if x[i] != y[i] {
			return false
		}
	}

	// test excess digits: nonzero => x != y
	for i := yz; i < xz; i++ {
		if x[i] != 0 {
			return false
		}
	}

	return true
}

// NotEqual is true if and only if x and y are not equal.
func NotEqual[Word constraints.Unsigned](x, y []Word) bool {
	// operation is commutative
	// ensure x is not shorter than y
	if len(x) < len(y) {
		x, y = y, x
	}
	xz := len(x)
	yz := len(y)
	// invariant: xz >= yz

	// test common digits, order doesn't matter
	for i := 0; i < yz; i++ {
		if x[i] != y[i] {
			return true
		}
	}

	// test excess digits: nonzero => x != y
	for i := yz; i < xz; i++ {
		if x[i] != 0 {
			return true
		}
	}

	return false
}
