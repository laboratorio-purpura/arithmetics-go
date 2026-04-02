// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "golang.org/x/exp/constraints"

// AreEqual tests if two integers are equal.
func AreEqual[Word constraints.Unsigned](x, y []Word) bool {
	xz := len(x)
	yz := len(y)

	z := min(xz, yz)

	for i := 0; i < z; i++ {
		if x[i] != y[i] {
			return false
		}
	}

	for i := z; i < xz; i++ {
		if x[i] != 0 {
			return false
		}
	}

	for i := z; i < yz; i++ {
		if y[i] != 0 {
			return false
		}
	}

	return true
}

// NotEqual tests if two integers are not equal.
func NotEqual[Word constraints.Unsigned](x, y []Word) bool {
	return !AreEqual(x, y)
}
