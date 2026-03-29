// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "golang.org/x/exp/constraints"

// IsZero returns true if and only if x is zero.
func IsZero[Word constraints.Unsigned](x []Word) bool {
	xz := len(x)
	for i := 0; i < xz; i++ {
		if x[i] != 0 {
			return false
		}
	}
	return true
}

// NotZero returns true if and only if x is nonzero.
func NotZero[Word constraints.Unsigned](x []Word) bool {
	xz := len(x)
	for i := 0; i < xz; i++ {
		if x[i] != 0 {
			return true
		}
	}
	return false
}
