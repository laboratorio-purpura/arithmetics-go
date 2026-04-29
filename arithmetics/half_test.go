// SPDX-FileCopyrightText: 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"math/bits"
	"testing"

	"pgregory.net/rapid"
)

func TestHalve_Differential_Rapid(t *testing.T) {
	const Bits = bits.UintSize

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
}

func BenchmarkHalve(b *testing.B) {
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
