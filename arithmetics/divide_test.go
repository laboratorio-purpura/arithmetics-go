// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math/big"
	"math/bits"
	"slices"
	"testing"

	"pgregory.net/rapid"
	"purpura.dev.br/arithmetics/arithmetics/internal"
)

func TestDivideNormalStrict2By1_Differential_Rapid(t *testing.T) {
	const Bits = bits.UintSize

	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		y := rapid.UintMin(1<<(Bits-1)).Draw(t, "y")
		isStrict := func(x []uint) bool {
			return x[1] < y
		}
		x := rapid.SliceOfN(rapid.Uint(), 2, 2).Filter(isStrict).Draw(t, "x")

		// compute with purple
		iy := Reciprocal(y)
		q, r := divideNormalStrict2By1([2]uint(x), y, iy)
		t.Logf("q = %v, r = %v", q, r)

		// compute with math/bits
		q_, r_ := bits.Div(x[1], x[0], y)
		t.Logf("q_ = %v, r_ = %v", q_, r_)

		// compare
		if q != q_ {
			t.Errorf("q = %v, q_ = %v", q, q_)
		}
		if r != r_ {
			t.Errorf("r = %v, r_ = %v", r, r_)
		}
	})
}

func TestDivideBy1_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		y := rapid.UintMin(1).Draw(t, "y1")
		t.Logf("y = %X", y)
		x := rapid.SliceOfN(rapid.Uint(), 1, 32).Draw(t, "x")
		t.Logf("x = %X", x)

		// compute with purple
		q := make([]uint, len(x))
		r := DivideBy1(q, x, y)
		t.Logf("q = %X, r = %X", q, r)

		// compute with math/big
		x_ := internal.ToBigInt(x)
		y_ := big.NewInt(0).SetUint64(uint64(y))
		q_ := big.NewInt(0).Div(x_, y_)
		r_ := big.NewInt(0).Mod(x_, y_)
		t.Logf("q_ = %X, r_ = %X", q_, r_)

		// compare
		if internal.ToBigInt(q).Cmp(q_) != 0 {
			t.Error("difference in quotient")
		}
		if big.NewInt(0).SetUint64(uint64(r)).Cmp(r_) != 0 {
			t.Error("difference in remainder")
		}
	})
}

func TestDivideNormalStrict3By2_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		y := [2]uint(rapid.SliceOfN(rapid.Uint(), 2, 2).Filter(IsNormal).Draw(t, "y"))
		t.Logf("y = %X", y)
		isStrict := func(i []uint) bool {
			return IsSmaller(i[1:3], y[:])
		}
		x := [3]uint(rapid.SliceOfN(rapid.Uint(), 3, 3).Filter(isStrict).Draw(t, "x"))
		t.Logf("x = %X", x)

		// compute with purple
		iy := Reciprocal2(y)
		q, r := divideNormalStrict3By2(x, y, iy)
		t.Logf("q = %X, r = %X", q, r)

		// compute with math/big
		x_ := internal.ToBigInt(x[:])
		y_ := internal.ToBigInt(y[:])
		q_ := big.NewInt(0).Div(x_, y_)
		r_ := big.NewInt(0).Mod(x_, y_)
		t.Logf("q_ = %X, r_ = %X", q_, r_)

		// compare
		if big.NewInt(0).SetUint64(uint64(q)).Cmp(q_) != 0 {
			t.Error("difference in quotient")
		}
		if internal.ToBigInt(r[:]).Cmp(r_) != 0 {
			t.Error("difference in remainder")
		}
	})
}

func TestDivideNormalStrictN1ByN(t *testing.T) {
	cases := []struct {
		x []uint
		y []uint
	}{
		{
			x: []uint{0x0, 0x0, 0xFFFFFFFFFFFFFFF9, 0xFFFFFFFFFFFFFFFB, 0x5, 0xFFFFFFFFFFFFFFFA},
			y: []uint{0x1, 0x1, 0x0, 0x1, 0xffffffffffffffff},
		},
	}
	for _, it := range cases {
		t.Run(fmt.Sprintf("%+v", it), func(t *testing.T) {
			t.Logf("x = %X", it.x)
			t.Logf("y = %X", it.y)

			iy := Reciprocal(it.y[len(it.y)-1])
			r := make([]uint, len(it.y))
			q := divideNormalStrictN1ByN(r, it.x, it.y, iy)
			t.Logf("q = %X", q)
			t.Logf("r = %X", r)

			x_ := internal.ToBigInt(it.x)
			y_ := internal.ToBigInt(it.y)
			q_ := big.NewInt(0).Div(x_, y_)
			r_ := big.NewInt(0).Mod(x_, y_)
			t.Logf("q_ = %X", q_)
			t.Logf("r_ = %X", r_)

			if big.NewInt(0).SetUint64(uint64(q)).Cmp(q_) != 0 {
				t.Error("difference in quotient")
			}
			if internal.ToBigInt(r).Cmp(r_) != 0 {
				t.Error("difference in remainder")
			}
		})
	}
}

func TestDivideNormalStrictN1ByN_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		N := rapid.IntRange(2, 32).Draw(t, "N")
		y := rapid.SliceOfN(rapid.Uint(), N, N).Filter(IsNormal).Draw(t, "y")
		isStrict := func(i []uint) bool {
			return IsSmaller(i[1:], y[:])
		}
		x := rapid.SliceOfN(rapid.Uint(), N+1, N+1).Filter(isStrict).Draw(t, "x")

		// compute with purple
		iy := Reciprocal(y[len(y)-1])
		r := make([]uint, len(x))
		q := divideNormalStrictN1ByN(r, x, y, iy)
		t.Logf("q = %X, r = %X", q, r)

		// compute with math/big
		x_ := internal.ToBigInt(x)
		y_ := internal.ToBigInt(y)
		q_ := big.NewInt(0).Div(x_, y_)
		r_ := big.NewInt(0).Mod(x_, y_)
		t.Logf("q_ = %X, r_ = %X", q_, r_)

		// compare
		if big.NewInt(0).SetUint64(uint64(q)).Cmp(q_) != 0 {
			t.Error("difference in quotient")
		}
		if internal.ToBigInt(r).Cmp(r_) != 0 {
			t.Error("difference in remainder")
		}
	})
}

func TestDivideNormalStrictN1ByN_Accumulate_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		N := rapid.IntRange(2, 32).Draw(t, "N")
		t.Logf("N = %v", N)
		y := rapid.SliceOfN(rapid.Uint(), N, N).Filter(IsNormal).Draw(t, "y")
		t.Logf("y = %v", y)
		isStrict := func(i []uint) bool {
			return IsSmaller(i[1:], y[:])
		}
		x := rapid.SliceOfN(rapid.Uint(), N+1, N+1).Filter(isStrict).Draw(t, "x")
		t.Logf("x = %v", x)

		iy := Reciprocal(y[len(y)-1])

		// compute in result style
		r1 := make([]uint, len(x))
		q1 := divideNormalStrictN1ByN(r1, x, y, iy)
		t.Logf("q = %v, r = %v", q1, r1)

		// compute in accumulation style
		r2 := make([]uint, len(x))
		copy(r2, x)
		q2 := divideNormalStrictN1ByN(r2, r2, y, iy)
		t.Logf("q2 = %v, r2 = %v", q2, r2)

		// compare
		if q1 != q2 {
			t.Error("difference in quotient")
		}
		if !slices.Equal(r1, r2) {
			t.Error("difference in remainder")
		}
	})
}

func TestDivide_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// generate samples
		y := rapid.SliceOfN(rapid.Uint(), 2, 32).Filter(IsCompact).Filter(NotZero).Draw(t, "y")
		x := rapid.SliceOfN(rapid.Uint(), 0, 32).Draw(t, "x")

		// compute with purple
		q := make([]uint, len(x))
		r := make([]uint, len(y))
		Divide(q, r, x, y)
		t.Logf("q = %X, r = %X", q, r)

		// compute with math/big
		x_ := internal.ToBigInt(x)
		y_ := internal.ToBigInt(y)
		t.Logf("x_ = %X, y_ = %X", x_, y_)
		q_ := big.NewInt(0).Div(x_, y_)
		r_ := big.NewInt(0).Mod(x_, y_)
		t.Logf("q_ = %X, r_ = %X", q_, r_)

		// compare
		if internal.ToBigInt(q).Cmp(q_) != 0 {
			t.Error("difference in quotient")
		}
		if internal.ToBigInt(r).Cmp(r_) != 0 {
			t.Error("difference in remainder")
		}
	})
}
