// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"testing"

	"pgregory.net/rapid"
)

func TestIsLess_ZeroNotLessThanZero_Rapid(t *testing.T) {
	// test zero is not less than zero
	rapid.Check(t, func(t *rapid.T) {
		words := rapid.IntRange(0, 64).Draw(t, "words")
		x := rapid.SliceOfN(rapid.UintMax(0), words, words).Draw(t, "x")
		y := rapid.SliceOfN(rapid.UintMax(0), words, words).Draw(t, "y")
		r := IsLess(x, y)
		if r {
			t.Error("IsLess(x,y) = true")
		}
	})
}

func TestIsLess_ZeroLessThanNonzero_Rapid(t *testing.T) {
	// test zero is less than nonzero
	rapid.Check(t, func(t *rapid.T) {
		words := rapid.IntRange(1, 64).Draw(t, "words")
		x := rapid.SliceOfN(rapid.UintMax(0), words, words).Draw(t, "x")
		y := rapid.SliceOfN(rapid.UintMin(1), words, words).Draw(t, "y")
		r := IsLess(x, y)
		if !r {
			t.Error("IsLess(x,y) = false")
		}
	})
}
