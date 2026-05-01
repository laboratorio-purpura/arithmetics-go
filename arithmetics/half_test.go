// SPDX-FileCopyrightText: 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"math/bits"
	"slices"
	"testing"

	"hegel.dev/go/hegel"
	"pgregory.net/rapid"
)

func TestHalfHegel(t *testing.T) {
	const Bits = bits.UintSize
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger())
			y := hegel.Draw[uint](ht, hegel.Integers[uint](0, Bits-1))
			z := len(x)
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			q := make([]uint, z)
			_ = Half(q, x, y)

			// compute with math/big

			x_ := toBigInt(x)
			q_ := big.NewInt(0).Rsh(x_, y)

			// compare

			if toBigInt(q).Cmp(q_) != 0 {
				ht.Fatalf("q = %X, q_ = %X", q, q_)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("accumulate", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger())
			y := hegel.Draw[uint](ht, hegel.Integers[uint](0, Bits-1))
			z := len(x)
			ht.Logf("x = %X, y = %X", x, y)

			// compute result

			q1 := make([]uint, z)
			r1 := Half(q1, x, y)

			// accumulate result

			q2 := make([]uint, z)
			copy(q2, x)
			r2 := Half(q2, q2, y)

			// compare

			if !slices.Equal(q1, q2) {
				ht.Fatalf("q1 = %X, q2 = %X", q1, q2)
			}
			if r1 != r2 {
				ht.Fatalf("r1 = %d, r2 = %d", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("short-result", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelNonemptyLongInteger())
			y := hegel.Draw[uint](ht, hegel.Integers[uint](0, Bits-1))
			// full size
			fz := len(x)
			// short size
			sz := hegel.Draw(ht, hegel.Integers[int](0, fz-1))
			ht.Logf("x = %X, y = %X, sz = %d", x, y, sz)

			// full result

			q1 := make([]uint, fz)
			_ = Half(q1, x, y)

			// short result

			q2 := make([]uint, sz)
			_ = Half(q2, x, y)

			// compare

			if !slices.Equal(q1[:sz], q2) {
				ht.Fatalf("q1 = %X, q2 = %X", q1, q2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestHalfRapid(t *testing.T) {
	const Bits = bits.UintSize
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
			t.Logf("x: %X", x)
			y := rapid.UintMax(Bits-1).Draw(t, "y")
			t.Logf("y: %X", y)

			// compute with purple
			rz := len(x)
			r := make([]uint, rz)
			Half(r, x, y)
			t.Logf("r: %X", r)

			// translate samples to math/big
			x_ := big.NewInt(0)
			for i := len(x); i > 0; i-- {
				x_.Lsh(x_, Bits)
				x_.Add(x_, big.NewInt(0).SetUint64(uint64(x[i-1])))
			}
			t.Logf("x_ = %X", x_)

			// compute with math/big
			r_ := big.NewInt(0).Rsh(x_, y)

			// compare
			for i, v := range r_.Bits() {
				if r[i] != uint(v) {
					t.Errorf("i = %d, r = %X, r_ = %X, r[i] = %X, r_[i] = %X", i, r, r_, r[i], v)
					break
				}
			}
		})
	})
}

func BenchmarkHalf(b *testing.B) {
	rng := newRand()

	for _, words := range []uint{8, 16, 32, 64, 128, 256} {
		// generate samples
		x := make([]uint, words)
		for i := range x {
			x[i] = rng.Uint()
		}
		y := rng.UintN(uint(bits.UintSize - 1))

		// measure purple
		b.Run(fmt.Sprint("purple-", words), func(b *testing.B) {
			t := make([]uint, words)
			for b.Loop() {
				Half(t, x, y)
			}
		})

		// measure math/big
		b.Run(fmt.Sprint("math-big-", words), func(b *testing.B) {
			x_ := toBigInt(x)
			t := big.NewInt(0)
			for b.Loop() {
				t.Rsh(x_, y)
			}
		})
	}
}
