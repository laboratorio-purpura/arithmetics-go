// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
	"testing"

	"pgregory.net/rapid"
)

func TestDivide2By1WithReciprocal32_Differential_Rapid(t *testing.T) {
	const Bits = bits.UintSize

	rapid.Check(t, func(t *rapid.T) {
		y := rapid.UintMin(1<<(Bits-1)).Draw(t, "y")
		x0 := rapid.Uint().Draw(t, "x0")
		x1 := rapid.UintMax(y-1).Draw(t, "x1")
		x := [2]uint{x0, x1}

		iy := Reciprocal(y)
		q, r := Divide2By1WithReciprocal(x, y, iy)
		t.Logf("q = %v, r = %v", q, r)

		q_, r_ := bits.Div(x1, x0, y)
		if q != q_ {
			t.Errorf("q = %v, q_ = %v", q, q_)
		}
		if r != r_ {
			t.Errorf("r = %v, r_ = %v", r, r_)
		}
	})
}

func TestDivideBy1WithReciprocal32_Definition_Rapid(t *testing.T) {
	const Bits = bits.UintSize

	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOfN(rapid.Uint(), 2, -1).Draw(t, "x")
		x[len(x)-1] |= 1 << (Bits - 1)
		t.Logf("x = %v", x)
		y := rapid.UintMin(1<<(Bits-1)).Draw(t, "y1")
		t.Logf("y = %v", y)

		iy := Reciprocal(y)

		q := make([]uint, len(x))
		r := DivideBy1WithReciprocal(q, x, y, iy)
		t.Logf("q = %v, r = %v", q, r)

		x_ := make([]uint, len(x)+1)
		Multiply(x_, []uint{y}, q)
		Add(x_, x_, []uint{r})
		if NotEqual(x_, x[:]) {
			t.Errorf("x_ = %v", x_)
		}
	})
}

func TestDivide3By2WithReciprocal32_Definition_Rapid(t *testing.T) {
	const Bits = bits.UintSize

	rapid.Check(t, func(t *rapid.T) {
		y0 := rapid.Uint().Draw(t, "y0")
		y1 := rapid.UintMin(1<<(Bits-1)).Draw(t, "y1")
		y := [2]uint{y0, y1}
		x0 := rapid.Uint().Draw(t, "x0")
		x1 := rapid.Uint().Draw(t, "x1")
		x2 := rapid.UintMax(y1-1).Draw(t, "x2")
		x := [3]uint{x0, x1, x2}

		iy := Reciprocal2(y)
		q, r := Divide3By2WithReciprocal(x, y, iy)
		t.Logf("q = %v, r = %v", q, r)

		x_ := make([]uint, 3)
		Multiply(x_, y[:], []uint{q})
		Add(x_, x_, r[:])
		if NotEqual(x_, x[:]) {
			t.Errorf("x_ = %v", x_)
		}
	})
}
