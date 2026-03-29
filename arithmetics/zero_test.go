// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"testing"

	"pgregory.net/rapid"
)

func TestIsZero(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		zero := rapid.SliceOf(rapid.UintMax(0)).Draw(t, "zero")
		r := IsZero(zero)
		if !r {
			t.Errorf("IsZero(%v) = true, want false", zero)
		}
	})
	rapid.Check(t, func(t *rapid.T) {
		nonzero := rapid.SliceOfN(rapid.UintMin(1), 1, -1).Draw(t, "nonzero")
		r := IsZero(nonzero)
		if r {
			t.Errorf("IsZero(%v) = true, want false", nonzero)
		}
	})
}

func TestNotZero(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		zero := rapid.SliceOf(rapid.UintMax(0)).Draw(t, "zero")
		r := NotZero(zero)
		if r {
			t.Errorf("NotZero(%v) = false, want true", zero)
		}
	})
	rapid.Check(t, func(t *rapid.T) {
		nonzero := rapid.SliceOfN(rapid.UintMin(1), 1, -1).Draw(t, "nonzero")
		r := NotZero(nonzero)
		if !r {
			t.Errorf("NotZero(%v) = true, want false", nonzero)
		}
	})
}
