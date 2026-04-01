// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"slices"
	"testing"

	"pgregory.net/rapid"
)

func TestProduct32_Commutativity_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")
		rz := len(x) + len(y)
		r1 := make([]uint, rz)
		Product(r1, x, y)
		r2 := make([]uint, rz)
		Product(r2, y, x)
		if !slices.Equal(r1, r2) {
			t.Error("Product(x,y) != Product(y,x)")
		}
	})
}

func TestProduct32_Identity_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOfN(rapid.Uint(), 1, -1).Draw(t, "x")
		identity := []uint{rapid.Just[uint](1).Draw(t, "identity")}
		rz := len(x) + len(identity)
		r1 := make([]uint, rz)
		Product(r1, x, identity)
		if !AreEqual(r1, x) {
			t.Error("Product(x,identity) != x")
		}
	})
}

func TestProduct32_Nihil_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		nihil := rapid.SliceOf(rapid.Just[uint](0)).Draw(t, "nihil")
		rz := len(x) + len(nihil)
		r := make([]uint, rz)
		Product(r, x, nihil)
		if !IsZero(r) {
			t.Error("Product(x,nihil) != nihil")
		}
	})
}

func TestProduct32_ResultLessThanPart_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOfN(rapid.UintMin(1), 1, -1).Draw(t, "x")
		y := rapid.SliceOfN(rapid.UintMin(1), 1, -1).Draw(t, "y")
		rz := len(x) + len(y)
		r := make([]uint, rz)
		Product(r, x, y)
		if IsLess(r, x) {
			t.Error("Product(x,y) < x")
		}
		if IsLess(r, x) {
			t.Error("Product(x,y) < y")
		}
	})
}
