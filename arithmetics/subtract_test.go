// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"slices"
	"testing"

	"pgregory.net/rapid"
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
		borrow := Difference(difference, x, y)
		t.Logf("difference = %X, borrow = %X", difference, borrow)

		// compute with math/big
		x_ := toBigInt(x)
		y_ := toBigInt(y)
		difference_ := big.NewInt(0).Sub(x_, y_)
		t.Logf("difference_ = %X", difference_)

		// compare
		if toBigInt(difference).Cmp(difference_) != 0 {
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
		borrow1 := Difference(result, x, y)
		t.Logf("difference = %X, borrow = %X", result, borrow1)

		// compute in accumulate style
		accumulator := make([]uint, max(len(x), len(y)))
		copy(accumulator, x)
		borrow2 := Difference(accumulator, accumulator, y)

		// compare
		if !slices.Equal(result, accumulator) {
			t.Error("difference in result")
		}
		if borrow1 != borrow2 {
			t.Error("difference in borrow")
		}
	})
}

func BenchmarkSubtract(b *testing.B) {
	rng := newRand()

	for _, words := range []uint{8, 16, 32, 64, 128, 256} {
		// generate samples
		x := make([]uint, words)
		for i := range x {
			x[i] = rng.Uint()
		}
		y := make([]uint, words)
		for i := range y {
			y[i] = rng.Uint()
		}

		// measure purple
		b.Run(fmt.Sprint("purple-", words), func(b *testing.B) {
			difference := make([]uint, words)
			var borrow uint
			for b.Loop() {
				borrow = Difference(difference, x, y)
			}
			_, _ = difference, borrow
		})

		// measure math/big
		b.Run(fmt.Sprint("math-big-", words), func(b *testing.B) {
			x_ := toBigInt(x)
			y_ := toBigInt(y)
			difference := big.NewInt(0)
			for b.Loop() {
				difference = difference.Sub(x_, y_)
			}
			_ = difference
		})
	}
}
