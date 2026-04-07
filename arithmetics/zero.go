// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

// IsZero tests if an integer is zero.
func IsZero(x []uint) bool {
	xz := len(x)

	// integer is zero if and only if every word is zero
	for i := 0; i < xz; i++ {
		if x[i] != 0 {
			return false
		}
	}

	return true
}

// NotZero tests if an integer is not zero.
func NotZero(x []uint) bool {
	xz := len(x)

	// integer is not zero if and only if any word is nonzero
	for i := 0; i < xz; i++ {
		if x[i] != 0 {
			return true
		}
	}

	return false
}
