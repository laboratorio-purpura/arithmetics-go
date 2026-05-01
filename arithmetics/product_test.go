// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math"
	"math/big"
	"slices"
	"testing"

	"hegel.dev/go/hegel"
	"pgregory.net/rapid"
)

func TestProductBy1Hegel(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger())
			y := hegel.Draw[uint](ht, hegel.Integers[uint](0, math.MaxUint))
			z := len(x)
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			r := make([]uint, z+1)
			r[z] = ProductBy1(r, x, y)

			// compute with math/big

			x_ := toBigInt(x)
			y_ := big.NewInt(0).SetUint64(uint64(y))
			r_ := big.NewInt(0).Mul(x_, y_)

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
			y := hegel.Draw[uint](ht, hegel.Integers[uint](0, math.MaxUint))
			z := len(x)
			ht.Logf("x = %X, y = %X", x, y)

			// compute result

			r1 := make([]uint, z)
			e1 := ProductBy1(r1, x, y)

			// accumulate result

			r2 := make([]uint, z)
			copy(r2, x)
			e2 := ProductBy1(r2, r2, y)

			// compare

			if !slices.Equal(r1, r2) {
				ht.Fatalf("r = %X, r_ = %X", r1, r2)
			}
			if e1 != e2 {
				ht.Fatalf("c1 = %d, c2 = %d", e1, e2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("short-result", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelNonemptyLongInteger())
			y := hegel.Draw[uint](ht, hegel.Integers[uint](0, math.MaxUint))
			// full size
			fz := len(x)
			// short size
			sz := hegel.Draw(ht, hegel.Integers[int](0, fz-1))
			ht.Logf("x = %X, y = %X, sz = %d", x, y, sz)

			// full result

			r1 := make([]uint, fz)
			_ = ProductBy1(r1, x, y)

			// short result

			r2 := make([]uint, sz)
			_ = ProductBy1(r2, x, y)

			// compare

			if !slices.Equal(r1[:sz], r2) {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestProductBy1Rapid(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
			y := rapid.Uint().Draw(t, "y")

			// compute with purple
			product := make([]uint, len(x)+1)
			product[len(x)] = ProductBy1(product, x, y)
			t.Logf("product = %X", product)

			// compute with math/big
			x_ := toBigInt(x)
			y_ := big.NewInt(0).SetUint64(uint64(y))
			product_ := big.NewInt(0).Mul(x_, y_)
			t.Logf("product_ = %X", product_)

			// compare
			if toBigInt(product).Cmp(product_) != 0 {
				t.Error("difference in product")
			}
		})
	})
}

func TestProductHegel(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger())
			y := hegel.Draw[[]uint](ht, hegelLongInteger())
			z := len(x) + len(y)
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			r := make([]uint, z)
			Product(r, x, y)

			// compute with math/big

			x_ := toBigInt(x)
			y_ := toBigInt(y)
			r_ := big.NewInt(0).Mul(x_, y_)

			// compare

			if toBigInt(r).Cmp(r_) != 0 {
				ht.Fatalf("r = %X, r_ = %X", r, r_)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("short-result", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelNonemptyLongInteger())
			y := hegel.Draw[[]uint](ht, hegelNonemptyLongInteger())
			// full size
			fz := len(x) + len(y)
			// short size
			sz := hegel.Draw(ht, hegel.Integers[int](0, fz-1))
			ht.Logf("x = %X, y = %X, sz = %d", x, y, sz)

			// full result

			r1 := make([]uint, fz)
			Product(r1, x, y)

			// short result

			r2 := make([]uint, sz)
			Product(r2, x, y)

			// compare

			if !slices.Equal(r1[:sz], r2) {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestProductRapid(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
			y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")

			// compute with purple
			product := make([]uint, len(x)+len(y))
			Product(product, x, y)
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
	})
}

func BenchmarkProduct(b *testing.B) {
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
				Product(product, x, y)
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
