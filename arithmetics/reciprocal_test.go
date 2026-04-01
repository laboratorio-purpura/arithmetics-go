package arithmetics

import (
	"math/big"
	"math/bits"
	"testing"

	"pgregory.net/rapid"
)

func TestReciprocalWord32_Differential_Rapid(t *testing.T) {
	const Bits = bits.UintSize

	rapid.Check(t, func(t *rapid.T) {
		y := rapid.UintMin(1<<(Bits-1)).Draw(t, "y")
		t.Logf("y = %X", y)

		iy := Reciprocal(y)
		t.Logf("iy = %X", iy)

		y_ := big.NewInt(0).SetUint64(uint64(y))
		t.Logf("y_ = %X", y_)
		if y_.Uint64() != uint64(y) {
			t.Fatal("y_ ≠ y")
		}

		// by definition: ((β^2 − 1) ÷ y) − β
		one := big.NewInt(1)
		β := big.NewInt(0).Lsh(one, Bits)
		β2 := big.NewInt(0).Lsh(one, Bits*2)
		iy_ := big.NewInt(0)
		iy_ = iy_.Sub(β2, one)
		iy_ = iy_.Div(iy_, y_)
		iy_ = iy_.Sub(iy_, β)
		t.Logf("iy_ = %X", iy_)
		if iy_.Cmp(β) != -1 {
			t.Fatal("iy_ ≥ β")
		}

		if iy_.Uint64() != uint64(iy) {
			t.Errorf("iy_ ≠ iy")
		}
	})
}

func TestReciprocalWords32_Differential_Rapid(t *testing.T) {
	const Bits = bits.UintSize

	rapid.Check(t, func(t *rapid.T) {
		y0 := rapid.Uint().Draw(t, "y0")
		y1 := rapid.UintMin(1<<(Bits-1)).Draw(t, "y1")
		y := [2]uint{y0, y1}
		t.Logf("y = %X", y)

		iy := Reciprocal2(y)
		t.Logf("iy = %X", iy)

		// y_ = y0 | y1 << Bits
		y_ := big.NewInt(0).SetUint64(uint64(y0))
		y_ = y_.Add(
			y_,
			big.NewInt(0).Lsh(
				big.NewInt(0).SetUint64(uint64(y1)),
				Bits,
			),
		)
		t.Logf("y_ = %X", y_)

		// by definition: iy_ ((β^3 − 1) ÷ y) − β
		_1 := big.NewInt(1)
		β := big.NewInt(0).Lsh(_1, Bits)
		β3 := big.NewInt(0).Lsh(_1, Bits*3)
		iy_ := big.NewInt(0)
		iy_ = iy_.Sub(β3, _1)
		iy_ = iy_.Div(iy_, y_)
		iy_ = iy_.Sub(iy_, β)
		t.Logf("iy_ = %X", iy_)
		if iy_.Cmp(β) != -1 {
			t.Fatal("iy_ ≥ β")
		}

		if iy_.Uint64() != uint64(iy) {
			t.Error("iy_ ≠ iy")
		}
	})
}
