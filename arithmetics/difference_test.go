// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"testing"

	"pgregory.net/rapid"
)

func TestDifference32_Identity_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		identity := rapid.SliceOf(rapid.Just[uint](0)).Draw(t, "identity")
		rz := max(len(x), len(identity))
		r := make([]uint, rz)
		b := Difference(r, x, identity)
		if !AreEqual(r, x) {
			t.Errorf("Difference(x,identity) != x")
		}
		if b != 0 {
			t.Errorf("Difference(x,identity) borrow != 0")
		}
	})
}

func TestDifference32_ResultLessThanPart_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")
		rz := len(x) + len(y)
		r := make([]uint, rz)
		b := Difference(r, x, y)
		if b == 0 {
			if IsGreater(r, x) {
				t.Error("x >= y but Difference(x,y) > x")
			}
		} else {
			if IsSmaller(r, x) {
				t.Error("x < y but Difference(x,y) < x")
			}
		}
	})
}
