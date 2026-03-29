// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"testing"

	"pgregory.net/rapid"
)

func TestEquals(t *testing.T) {
	samples := []struct {
		x []uint
		y []uint
		r bool
	}{
		{[]uint{1}, []uint{}, false},
	}

	for _, sample := range samples {
		t.Run(fmt.Sprintf("%+v", sample), func(t *testing.T) {
			r := AreEqual(sample.x, sample.y)
			if r != sample.r {
				t.Fatal("expected", sample.r, "got", r)
			}
		})
	}
}

func TestEqualsCommutativity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")
		r1 := AreEqual(x, y)
		r2 := AreEqual(y, x)
		if r1 != r2 {
			t.Fatal("Equal(x, y) != Equal(y, x)")
		}
	})
}

func TestNotEqualsCommutativity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")
		r1 := NotEqual(x, y)
		r2 := NotEqual(y, x)
		if r1 != r2 {
			t.Fatal("NotEqual(x, y) != NotEqual(y, x)")
		}
	})
}
