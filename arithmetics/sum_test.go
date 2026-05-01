// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"

	"hegel.dev/go/hegel"
	"pgregory.net/rapid"
)

func TestSumHegel(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger())
			y := hegel.Draw[[]uint](ht, hegelLongInteger())
			z := max(len(x), len(y))
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			r := make([]uint, z+1)
			r[z] = Sum(r, x, y)

			// compute with math/big

			x_ := toBigInt(x)
			y_ := toBigInt(y)
			r_ := big.NewInt(0).Add(x_, y_)

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

			// result

			r1 := make([]uint, z)
			c1 := Sum(r1, x, y)

			// accumulate result

			r2 := make([]uint, z)
			copy(r2, x)
			c2 := Sum(r2, r2, y)

			// compare

			if !slices.Equal(r1, r2) {
				ht.Fatalf("r = %X, r_ = %X", r1, r2)
			}
			if c1 != c2 {
				ht.Fatalf("c1 = %d, c2 = %d", c1, c2)
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
			_ = Sum(r1, x, y)

			// short result

			r2 := make([]uint, sz)
			_ = Sum(r2, x, y)

			// compare

			if !slices.Equal(r1[:sz], r2) {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestSumRapid(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
			y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")

			// compute with purple
			z := max(len(x), len(y)) + 1
			sum := make([]uint, z)
			sum[z-1] = Sum(sum, x, y)
			t.Logf("sum = %X", sum)

			// compute with math/big
			x_ := toBigInt(x)
			y_ := toBigInt(y)
			sum_ := big.NewInt(0).Add(x_, y_)
			t.Logf("sum_ = %X", sum_)

			// compare
			if toBigInt(sum).Cmp(sum_) != 0 {
				t.Error("difference in sum")
			}
		})
	})
	t.Run("accumulate", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
			y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")

			z := max(len(x), len(y))

			// compute in result style
			sum1 := make([]uint, z)
			carry1 := Sum(sum1, x, y)
			t.Logf("sum1 = %X, carry1 = %X", sum1, carry1)

			// compute in accumulate style
			sum2 := make([]uint, z)
			copy(sum2, x)
			carry2 := Sum(sum2, sum2, y)
			t.Logf("sum2 = %X, carry2 = %X", sum2, carry2)

			// compare
			if !slices.Equal(sum1, sum2) {
				t.Error("difference in sum")
			}
			if carry1 != carry2 {
				t.Error("difference in carry")
			}
		})
	})
}

func BenchmarkSum(b *testing.B) {
	for _, words := range []uint{8, 16, 32, 64, 128, 256} {
		// generate samples
		rng := rand.New(rand.NewPCG(31, 39))
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
			sum := make([]uint, words)
			var carry uint
			for b.Loop() {
				carry = Sum(sum, x, y)
			}
			_, _ = sum, carry
		})

		// measure math/big
		b.Run(fmt.Sprint("math-big-", words), func(b *testing.B) {
			x_ := toBigInt(x)
			y_ := toBigInt(y)
			sum := big.NewInt(0)
			for b.Loop() {
				sum.Add(x_, y_)
			}
			_ = sum
		})
	}
}
