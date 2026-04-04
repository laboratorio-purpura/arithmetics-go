// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/big"
	"slices"
	"testing"

	"pgregory.net/rapid"
	"purpura.dev.br/arithmetics/arithmetics/internal"
)

func TestSubtract_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		x := rapid.SliceOfN(rapid.Uint(), 1, -1).Draw(t, "x")
		isSmaller := func(i []uint) bool {
			return IsSmaller(i, x)
		}
		y := rapid.SliceOfN(rapid.Uint(), 1, -1).Filter(isSmaller).Draw(t, "y")
		// TODO: test negative results

		// compute with purple
		difference := make([]uint, max(len(x), len(y)))
		borrow := Subtract(difference, x, y)
		t.Logf("difference = %X, borrow = %X", difference, borrow)

		// compute with math/big
		x_ := internal.ToBigInt(x)
		y_ := internal.ToBigInt(y)
		difference_ := big.NewInt(0).Sub(x_, y_)
		t.Logf("difference_ = %X", difference_)

		// compare
		if internal.ToBigInt(difference).Cmp(difference_) != 0 {
			t.Error("difference")
		}
	})
}

func TestSubtract_Accumulate_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		x := rapid.SliceOfN(rapid.Uint(), 1, -1).Draw(t, "x")
		y := rapid.SliceOfN(rapid.Uint(), 1, -1).Draw(t, "y")

		// compute in result style
		result := make([]uint, max(len(x), len(y)))
		borrow1 := Subtract(result, x, y)
		t.Logf("difference = %X, borrow = %X", result, borrow1)

		// compute in accumulate style
		accumulator := make([]uint, max(len(x), len(y)))
		copy(accumulator, x)
		borrow2 := Subtract(accumulator, accumulator, y)

		// compare
		if !slices.Equal(result, accumulator) {
			t.Error("difference in result")
		}
		if borrow1 != borrow2 {
			t.Error("difference in borrow")
		}
	})
}
