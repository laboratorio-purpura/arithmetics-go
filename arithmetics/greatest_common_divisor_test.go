// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/big"
	"testing"

	"hegel.dev/go/hegel"
)

func TestGreatestCommonDivisorClobberHegel(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelNonemptyLongInteger())
			y := hegel.Draw[[]uint](ht, hegelNonemptyLongInteger())
			ht.Logf("x = %X, y = %X", x, y)

			// compute with math/big

			x_ := toBigInt(x)
			y_ := toBigInt(y)
			r_ := big.NewInt(0).GCD(nil, nil, x_, y_)

			// compute with purple — clobbers x and y

			r := GreatestCommonDivisorClobber(x, y)

			// compare

			if toBigInt(r).Cmp(r_) != 0 {
				ht.Fatalf("r = %X, r_ = %X", r, r_)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}
