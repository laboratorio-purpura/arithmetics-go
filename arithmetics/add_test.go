// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"

	"pgregory.net/rapid"
)

func TestAdd_Differential_Rapid(t *testing.T) {
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
}

func TestAdd_Accumulate_Rapid(t *testing.T) {
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
}

func BenchmarkAdd(b *testing.B) {
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
