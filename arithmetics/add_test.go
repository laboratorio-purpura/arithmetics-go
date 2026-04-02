// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"math/bits"
	"math/rand/v2"
	"slices"
	"testing"

	"pgregory.net/rapid"
)

func TestAdd_Commutativity_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")
		rz := len(x) + len(y)
		r1 := make([]uint, rz)
		c1 := Add(r1, x, y)
		r2 := make([]uint, rz)
		c2 := Add(r2, y, x)
		if !slices.Equal(r1, r2) {
			t.Error("Sum(x,y) != Sum(y,x)")
		}
		if c1 != c2 {
			t.Error("Sum(x,y) carry != Sum(y,x) carry")
		}
	})
}

func TestAdd_Identity_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		identity := rapid.SliceOf(rapid.Just[uint](0)).Draw(t, "identity")
		rz := len(x) + len(identity)
		r := make([]uint, rz)
		c := Add(r, x, identity)
		if !AreEqual(r, x) {
			t.Error("Sum(x,identity) != x")
		}
		if c != 0 {
			t.Error("Sum(x,identity) carry != 0")
		}
	})
}

func TestAdd_ResultLessThanPart_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOfN(rapid.UintMin(1), 1, -1).Draw(t, "x")
		y := rapid.SliceOfN(rapid.UintMin(1), 1, -1).Draw(t, "y")
		rz := len(x) + len(y) + 1
		r := make([]uint, rz)
		r[rz-1] = Add(r, x, y)
		if IsSmaller(r, x) {
			t.Error("Sum(x,y) < x")
		}
		if IsSmaller(r, y) {
			t.Error("Sum(x,y) < y")
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

		// translate samples to math/big
		x_ := big.NewInt(0)
		for i := len(x); i > 0; i-- {
			w := big.NewInt(0).SetUint64(uint64(x[i-1]))
			x_.Lsh(x_, bits.UintSize).Add(x_, w)
		}
		y_ := big.NewInt(0)
		for i := len(y); i > 0; i-- {
			w := big.NewInt(0).SetUint64(uint64(y[i-1]))
			y_.Lsh(y_, bits.UintSize).Add(y_, w)
		}

		// measure purple
		b.Run(fmt.Sprint("purple", words), func(b *testing.B) {
			t := make([]uint, words)
			for b.Loop() {
				Add(t, x, y)
			}
		})

		// measure math/big
		b.Run(fmt.Sprint("math-big-", words), func(b *testing.B) {
			t := big.NewInt(0)
			for b.Loop() {
				t.Add(x_, y_)
			}
		})
	}
}
