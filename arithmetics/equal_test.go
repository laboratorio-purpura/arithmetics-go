// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"testing"

	"pgregory.net/rapid"
)

func TestEqual_Rapid(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
			y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")

			// compute with purple
			equals := AreEqual(x, y)
			t.Logf("equals = %v", equals)

			// compute with math/big
			x_ := toBigInt(x)
			y_ := toBigInt(y)
			equals_ := x_.CmpAbs(y_) == 0
			t.Logf("equals_ = %v", equals_)

			// compare
			if equals != equals {
				t.Error("difference in result")
			}
		})
	})
}

func BenchmarkEqual(b *testing.B) {
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
			var equals bool
			for b.Loop() {
				equals = AreEqual(x, y)
			}
			_ = equals
		})

		// measure math/big
		b.Run(fmt.Sprint("math-big-", words), func(b *testing.B) {
			x_ := toBigInt(x)
			y_ := toBigInt(y)
			var equals bool
			for b.Loop() {
				equals = x_.CmpAbs(y_) == 0
			}
			_ = equals
		})
	}
}
