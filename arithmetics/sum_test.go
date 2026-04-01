// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"slices"
	"testing"

	"pgregory.net/rapid"
)

func TestSum_Commutativity_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")
		rz := len(x) + len(y)
		r1 := make([]uint, rz)
		c1 := Sum(r1, x, y)
		r2 := make([]uint, rz)
		c2 := Sum(r2, y, x)
		if !slices.Equal(r1, r2) {
			t.Error("Sum(x,y) != Sum(y,x)")
		}
		if c1 != c2 {
			t.Error("Sum(x,y) carry != Sum(y,x) carry")
		}
	})
}

func TestSum_Identity_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		identity := rapid.SliceOf(rapid.Just[uint](0)).Draw(t, "identity")
		rz := len(x) + len(identity)
		r := make([]uint, rz)
		c := Sum(r, x, identity)
		if !AreEqual(r, x) {
			t.Error("Sum(x,identity) != x")
		}
		if c != 0 {
			t.Error("Sum(x,identity) carry != 0")
		}
	})
}

func TestSum_ResultLessThanPart_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOfN(rapid.UintMin(1), 1, -1).Draw(t, "x")
		y := rapid.SliceOfN(rapid.UintMin(1), 1, -1).Draw(t, "y")
		rz := len(x) + len(y) + 1
		r := make([]uint, rz)
		r[rz-1] = Sum(r, x, y)
		if IsLess(r, x) {
			t.Error("Sum(x,y) < x")
		}
		if IsLess(r, y) {
			t.Error("Sum(x,y) < y")
		}
	})
}
