// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"slices"
	"testing"

	"hegel.dev/go/hegel"
	"pgregory.net/rapid"
)

func TestDifferenceHegel(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger())
			y := hegel.Draw[[]uint](ht, hegelLongInteger())
			z := max(len(x), len(y))
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			r := make([]uint, z)
			b := Difference(r, x, y)

			// TODO: lift this restriction
			ht.Assume(b == 0)

			// compute with math/big

			x_ := toBigInt(x)
			y_ := toBigInt(y)
			r_ := big.NewInt(0).Sub(x_, y_)

			// compare

			if toBigInt(r).Cmp(r_) != 0 {
				ht.Fatalf("r = %X, r_ = %X", r, r_)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("accumulate", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger())
			y := hegel.Draw[[]uint](ht, hegelLongInteger())
			z := max(len(x), len(y))
			ht.Logf("x = %X, y = %X", x, y)

			// compute result

			r1 := make([]uint, z)
			b1 := Difference(r1, x, y)

			// accumulate result

			r2 := make([]uint, z)
			copy(r2, x)
			b2 := Difference(r2, r2, y)

			// compare

			if !slices.Equal(r1, r2) {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}
			if b1 != b2 {
				ht.Fatalf("b1 = %d, b2 = %d", b1, b2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("short-result", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelNonemptyLongInteger())
			y := hegel.Draw[[]uint](ht, hegelNonemptyLongInteger())
			// full size
			fz := max(len(x), len(y))
			// short size
			sz := hegel.Draw(ht, hegel.Integers[int](0, fz-1))
			ht.Logf("x = %X, y = %X, sz = %d", x, y, sz)

			// full result

			r1 := make([]uint, fz)
			_ = Difference(r1, x, y)

			// short result

			r2 := make([]uint, sz)
			_ = Difference(r2, x, y)

			// compare

			if !slices.Equal(r1[:sz], r2) {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestDifferenceRapid(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
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
	})
	t.Run("accumulate", func(t *testing.T) {
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
	})
}

func BenchmarkDifference(b *testing.B) {
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
