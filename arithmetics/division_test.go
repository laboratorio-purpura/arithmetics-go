// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
	"testing"

	"pgregory.net/rapid"
)

func TestDivision2By1WithReciprocal32_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		y := rapid.Uint32Min(1<<31).Draw(t, "y")
		x0 := rapid.Uint32().Draw(t, "x0")
		x1 := rapid.Uint32Max(y-1).Draw(t, "x1")
		x := [2]uint32{x0, x1}

		iy := Reciprocal1W32(y)
		q, r := Division2By1WithReciprocal32(x, y, iy)
		t.Logf("q = %v, r = %v", q, r)

		q_, r_ := bits.Div32(x1, x0, y)
		if q != q_ {
			t.Errorf("q = %v, q_ = %v", q, q_)
		}
		if r != r_ {
			t.Errorf("r = %v, r_ = %v", r, r_)
		}
	})
}

func TestDivisionBy1WithReciprocal32_Definition_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOfN(rapid.Uint32(), 2, -1).Draw(t, "x")
		x[len(x)-1] |= 1 << 31
		t.Logf("x = %v", x)
		y := rapid.Uint32Min(1<<31).Draw(t, "y1")
		t.Logf("y = %v", y)

		iy := Reciprocal1W32(y)

		q := make([]uint32, len(x))
		r := DivisionBy1WithReciprocal32(q, x, y, iy)
		t.Logf("q = %v, r = %v", q, r)

		x_ := make([]uint32, len(x)+1)
		Product32(x_, []uint32{y}, q)
		Sum32(x_, x_, []uint32{r})
		if NotEqual(x_, x[:]) {
			t.Errorf("x_ = %v", x_)
		}
	})
}

func TestDivision3By2WithReciprocal32_Definition_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		y0 := rapid.Uint32().Draw(t, "y0")
		y1 := rapid.Uint32Min(1<<31).Draw(t, "y1")
		y := [2]uint32{y0, y1}
		x0 := rapid.Uint32().Draw(t, "x0")
		x1 := rapid.Uint32().Draw(t, "x1")
		x2 := rapid.Uint32Max(y1-1).Draw(t, "x2")
		x := [3]uint32{x0, x1, x2}

		iy := Reciprocal2W32(y)
		q, r := Division3By2WithReciprocal32(x, y, iy)
		t.Logf("q = %v, r = %v", q, r)

		x_ := make([]uint32, 3)
		Product32(x_, y[:], []uint32{q})
		Sum32(x_, x_, r[:])
		if NotEqual(x_, x[:]) {
			t.Errorf("x_ = %v", x_)
		}
	})
}
