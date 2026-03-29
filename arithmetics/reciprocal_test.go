package arithmetics

import (
	"math/big"
	"testing"

	"pgregory.net/rapid"
)

func TestReciprocalWord32_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		y := rapid.Uint32Min(1<<31).Draw(t, "y")
		t.Logf("y = %X", y)

		iy := Reciprocal1W32(y)
		t.Logf("iy = %X", iy)

		y_ := big.NewInt(0).SetUint64(uint64(y))
		t.Logf("y_ = %X", y_)
		if y_.Uint64() != uint64(y) {
			t.Fatal("y_ ≠ y")
		}

		// by definition: ((β^2 − 1) ÷ y) − β
		one := big.NewInt(1)
		β := big.NewInt(0).Lsh(one, 32)
		β2 := big.NewInt(0).Lsh(one, 32*2)
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

func TestReciprocalWords32(t *testing.T) {
	y := [2]uint32{0x13, 0x80000003}
	iy := Reciprocal2W32(y)
	if iy != 0xFFFFFFF3 {
		t.Errorf("expected FFFFFFF3, got %X", iy)
	}
}

func TestReciprocalWords32_Differential_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		y0 := rapid.Uint32().Draw(t, "y0")
		y1 := rapid.Uint32Min(1<<31).Draw(t, "y1")
		y := [2]uint32{y0, y1}
		t.Logf("y = %X", y)

		iy := Reciprocal2W32(y)
		t.Logf("iy = %X", iy)

		// y_ = y0 | y1 << 32
		y_ := big.NewInt(0).SetUint64(uint64(y0))
		y_ = y_.Add(
			y_,
			big.NewInt(0).Lsh(
				big.NewInt(0).SetUint64(uint64(y1)),
				32,
			),
		)
		t.Logf("y_ = %X", y_)

		// by definition: iy_ ((β^3 − 1) ÷ y) − β
		_1 := big.NewInt(1)
		β := big.NewInt(0).Lsh(_1, 32)
		β3 := big.NewInt(0).Lsh(_1, 32*3)
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
