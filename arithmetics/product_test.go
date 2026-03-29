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
		x := rapid.SliceOf(rapid.Uint32()).Draw(t, "x")
		y := rapid.SliceOf(rapid.Uint32()).Draw(t, "y")
		rz := len(x) + len(y)
		r1 := make([]uint32, rz)
		Product32(r1, x, y)
		r2 := make([]uint32, rz)
		Product32(r2, y, x)
		if !slices.Equal(r1, r2) {
			t.Error("Product(x,y) != Product(y,x)")
		}
	})
}

func TestProduct32_Identity_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOfN(rapid.Uint32(), 1, -1).Draw(t, "x")
		identity := []uint32{rapid.Just[uint32](1).Draw(t, "identity")}
		rz := len(x) + len(identity)
		r1 := make([]uint32, rz)
		Product32(r1, x, identity)
		if !AreEqual(r1, x) {
			t.Error("Product(x,identity) != x")
		}
	})
}

func TestProduct32_Nihil_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint32()).Draw(t, "x")
		nihil := rapid.SliceOf(rapid.Just[uint32](0)).Draw(t, "nihil")
		rz := len(x) + len(nihil)
		r := make([]uint32, rz)
		Product32(r, x, nihil)
		if !IsZero(r) {
			t.Error("Product(x,nihil) != nihil")
		}
	})
}

func TestProduct32_ResultLessThanPart_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOfN(rapid.Uint32Min(1), 1, -1).Draw(t, "x")
		y := rapid.SliceOfN(rapid.Uint32Min(1), 1, -1).Draw(t, "y")
		rz := len(x) + len(y)
		r := make([]uint32, rz)
		Product32(r, x, y)
		if IsLess(r, x) {
			t.Error("Product(x,y) < x")
		}
		if IsLess(r, x) {
			t.Error("Product(x,y) < y")
		}
	})
}
