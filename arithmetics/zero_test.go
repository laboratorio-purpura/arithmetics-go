// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"testing"

	"pgregory.net/rapid"
)

func TestIsZero_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")

		// compute with purple
		zero := IsZero(x)
		t.Logf("equals = %v", zero)

		// compute with math/big
		x_ := toBigInt(x)
		y_ := big.NewInt(0)
		zero_ := x_.CmpAbs(y_) == 0
		t.Logf("equals_ = %v", zero_)

		// compare
		if zero != zero {
			t.Error("difference in result")
		}
	})
}

func BenchmarkIsZero(b *testing.B) {
	rng := newRand()

	for _, words := range []uint{8, 16, 32, 64, 128, 256} {
		// generate samples
		x := make([]uint, words)
		for i := range x {
			x[i] = rng.Uint()
		}

		// translate samples to math/big

		// measure purple
		b.Run(fmt.Sprint("purple-", words), func(b *testing.B) {
			var zero bool
			for b.Loop() {
				zero = IsZero(x)
			}
			_ = zero
		})

		// measure math/big
		b.Run(fmt.Sprint("math-big-", words), func(b *testing.B) {
			x_ := toBigInt(x)
			y_ := big.NewInt(0)
			var zero bool
			for b.Loop() {
				zero = x_.CmpAbs(y_) == 0
			}
			_ = zero
		})
	}
}
