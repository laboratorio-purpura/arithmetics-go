// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/big"
	"math/bits"
	"testing"

	"pgregory.net/rapid"
)

func TestIsLess_Differential_Rapid(t *testing.T) {
	const Bits = bits.UintSize

	rapid.Check(t, func(t *rapid.T) {
		x := rapid.SliceOf(rapid.Uint()).Draw(t, "x")
		y := rapid.SliceOf(rapid.Uint()).Draw(t, "y")
		r := IsSmaller(x, y)
		t.Logf("r = %v", r)

		x_ := big.NewInt(0)
		for i := len(x); i > 0; i-- {
			x_.Lsh(x_, Bits)
			x_.Add(x_, big.NewInt(0).SetUint64(uint64(x[i-1])))
		}
		t.Logf("x_ = %v", x_)
		y_ := big.NewInt(0)
		for i := len(y); i > 0; i-- {
			y_.Lsh(y_, Bits)
			y_.Add(y_, big.NewInt(0).SetUint64(uint64(y[i-1])))
		}
		t.Logf("y_ = %v", y_)
		r_ := x_.Cmp(y_) < 0
		t.Logf("r_ = %v", r_)

		if r_ != r {
			t.Errorf("r_ != r")
		}
	})
}
