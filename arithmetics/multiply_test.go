// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"testing"

	"pgregory.net/rapid"
)

func TestMultiply_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")

		// compute with purple
		product := make([]uint, len(x)+len(y))
		Multiply(product, x, y)
		t.Logf("product = %X", product)

		// compute with math/big
		x_ := toBigInt(x)
		y_ := toBigInt(y)
		product_ := big.NewInt(0).Mul(x_, y_)
		t.Logf("product_ = %X", product_)

		// compare
		if toBigInt(product).Cmp(product_) != 0 {
			t.Error("difference in product")
		}
	})
}

func BenchmarkMultiply(b *testing.B) {
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
			product := make([]uint, words*2)
			for b.Loop() {
				Multiply(product, x, y)
			}
			_ = product
		})

		// measure math/big
		b.Run(fmt.Sprint("math-big-", words), func(b *testing.B) {
			x_ := toBigInt(x)
			y_ := toBigInt(y)
			product := big.NewInt(0)
			for b.Loop() {
				product = product.Mul(x_, y_)
			}
			_ = product
		})
	}
}
