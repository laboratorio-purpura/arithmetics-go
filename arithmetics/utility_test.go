// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math"
	"math/big"
	"math/bits"
	"math/rand/v2"
	"unsafe"

	"golang.org/x/exp/constraints"
	"hegel.dev/go/hegel"
	"pgregory.net/rapid"
)

const hegelCases = 10000

func hegelLongInteger() hegel.Generator[[]uint] {
	return hegel.Lists(hegel.Integers[uint](0, math.MaxUint))
}

func hegelNonemptyLongInteger() hegel.Generator[[]uint] {
	return hegel.Filter(hegel.Lists(hegel.Integers[uint](0, math.MaxUint)), func(v []uint) bool {
		return len(v) > 0
	})
}

func hegelNonzeroCompactLongInteger() hegel.Generator[[]uint] {
	return hegel.Filter(hegel.Lists(hegel.Integers[uint](0, math.MaxUint)), func(v []uint) bool {
		return NotZero(v) && IsCompact(v)
	})
}

func hegelNormalLongInteger() hegel.Generator[[]uint] {
	return hegel.Filter(hegel.Lists(hegel.Integers[uint](0, math.MaxUint)), func(v []uint) bool {
		return IsNormal(v)
	})
}

func normalLongInteger[E constraints.Integer](unit *rapid.Generator[E], minLen int, maxLen int) *rapid.Generator[[]E] {
	return rapid.Custom(func(t *rapid.T) []E {
		v := rapid.SliceOfN(unit, minLen, maxLen).Draw(t, "")
		v[len(v)-1] |= E(1) << (unsafe.Sizeof(E(0))*8 - 1)
		return v
	})
}

func newRand() *rand.Rand {
	return rand.New(rand.NewPCG(31, 39))
}

func toBigInt(x []uint) *big.Int {
	x_ := big.NewInt(0)
	for i := len(x); i > 0; i-- {
		w := big.NewInt(0).SetUint64(uint64(x[i-1]))
		x_.Lsh(x_, bits.UintSize).Add(x_, w)
	}
	return x_
}
