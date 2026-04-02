// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

// IsSmaller tests if an integer is smaller than another.
func IsSmaller(x, y []uint) bool {
	xz := len(x)
	yz := len(y)

	z := min(xz, yz)

	for i := xz; i > z; i-- {
		if x[i-1] != 0 {
			return false
		}
	}

	for i := yz; i > z; i-- {
		if y[i-1] != 0 {
			return true
		}
	}

	for i := z; i > 0; i-- {
		if x[i-1] < y[i-1] {
			return true
		}
		if x[i-1] > y[i-1] {
			return false
		}
	}

	return false
}

// NotSmaller tests if an integer is not smaller than another.
func NotSmaller(x, y []uint) bool {
	return !IsSmaller(x, y)
}
