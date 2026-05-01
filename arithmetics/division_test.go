// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"slices"
	"testing"

	"hegel.dev/go/hegel"
	"pgregory.net/rapid"
)

func TestDivisionNormalStrict2By1Hegel(t *testing.T) {
	const Bits = bits.UintSize
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			// y must be "normal"
			const normal = 1 << (Bits - 1)
			y := hegel.Draw(ht, hegel.Integers[uint](normal, math.MaxUint))
			// x must be "strict": x ÷ β < y
			strict := y - 1
			x := [2]uint{
				hegel.Draw(ht, hegel.Integers[uint](0, math.MaxUint)),
				hegel.Draw(ht, hegel.Integers[uint](0, strict)),
			}
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			iy := Reciprocal(y)
			q, r := divisionNormalStrict2By1(x, y, iy)

			// compute with math/bits

			q_, r_ := bits.Div(x[1], x[0], y)

			// compare

			if q != q_ {
				t.Errorf("q = %v, q_ = %v", q, q_)
			}
			if r != r_ {
				t.Errorf("r = %v, r_ = %v", r, r_)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestDivisionNormalStrict2By1Rapid(t *testing.T) {
	const Bits = bits.UintSize
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			y := rapid.UintMin(1<<(Bits-1)).Draw(t, "y")
			isStrict := func(x []uint) bool {
				return x[1] < y
			}
			x := rapid.SliceOfN(rapid.Uint(), 2, 2).Filter(isStrict).Draw(t, "x")

			// compute with purple
			iy := Reciprocal(y)
			q, r := divisionNormalStrict2By1([2]uint(x), y, iy)
			t.Logf("q = %v, r = %v", q, r)

			// compute with math/bits
			q_, r_ := bits.Div(x[1], x[0], y)
			t.Logf("q_ = %v, r_ = %v", q_, r_)

			// compare
			if q != q_ {
				t.Errorf("q = %v, q_ = %v", q, q_)
			}
			if r != r_ {
				t.Errorf("r = %v, r_ = %v", r, r_)
			}
		})
	})
}

func TestDivisionBy1Hegel(t *testing.T) {
	t.Run("accumulate", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger2(0, 64))
			y := hegel.Draw(ht, hegel.Integers[uint](1, math.MaxUint))
			ht.Logf("x = %X, y = %X", x, y)

			// compute result

			q1 := make([]uint, len(x))
			r1 := DivisionUni(q1, x, y)

			// accumulate result

			q2 := make([]uint, len(x))
			copy(q2, x)
			r2 := DivisionUni(q2, q2, y)

			// compare

			if NotEqual(q1, q2) {
				ht.Fatalf("q1 = %X, q2 = %X", q1, q2)
			}
			if r1 != r2 {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger2(0, 64))
			y := hegel.Draw(ht, hegel.Integers[uint](1, math.MaxUint))
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			q := make([]uint, len(x))
			r := DivisionUni(q, x, y)

			// compute with math/big

			x_ := toBigInt(x)
			y_ := big.NewInt(0).SetUint64(uint64(y))
			q_ := big.NewInt(0).Div(x_, y_)
			r_ := big.NewInt(0).Mod(x_, y_)

			// compare

			if toBigInt(q).Cmp(q_) != 0 {
				ht.Fatalf("q = %X, q_ = %X", q, q_)
			}
			if big.NewInt(0).SetUint64(uint64(r)).Cmp(r_) != 0 {
				ht.Fatalf("r = %X, r_ = %X", r, r_)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("short-quotient", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw(ht, hegelLongInteger2(1, 64))
			y := hegel.Draw(ht, hegel.Integers[uint](1, math.MaxUint))
			// full size
			fz := len(x)
			// short size
			sz := hegel.Draw(ht, hegel.Integers[int](0, fz-1))
			ht.Logf("x = %X, y = %X, sz = %d", x, y, sz)

			// full result

			q1 := make([]uint, fz)
			_ = DivisionUni(q1, x, y)

			// short result

			q2 := make([]uint, sz)
			_ = DivisionUni(q2, x, y)

			// compare

			if NotEqual(q1[:sz], q2) {
				ht.Fatalf("q1 = %X, q2 = %X", q1, q2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestDivisionBy1Rapid(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			y := rapid.UintMin(1).Draw(t, "y1")
			t.Logf("y = %X", y)
			x := rapid.SliceOfN(rapid.Uint(), 1, 32).Draw(t, "x")
			t.Logf("x = %X", x)

			// compute with purple
			q := make([]uint, len(x))
			r := DivisionUni(q, x, y)
			t.Logf("q = %X, r = %X", q, r)

			// compute with math/big
			x_ := toBigInt(x)
			y_ := big.NewInt(0).SetUint64(uint64(y))
			q_ := big.NewInt(0).Div(x_, y_)
			r_ := big.NewInt(0).Mod(x_, y_)
			t.Logf("q_ = %X, r_ = %X", q_, r_)

			// compare
			if toBigInt(q).Cmp(q_) != 0 {
				t.Error("difference in quotient")
			}
			if big.NewInt(0).SetUint64(uint64(r)).Cmp(r_) != 0 {
				t.Error("difference in remainder")
			}
		})
	})
}

func TestDivisionNormalStrict3By2Hegel(t *testing.T) {
	const Bits = bits.UintSize
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			// y must be "normal"
			const normal = 1 << (Bits - 1)
			y := [2]uint{
				hegel.Draw(ht, hegel.Integers[uint](0, math.MaxUint)),
				hegel.Draw(ht, hegel.Integers[uint](normal, math.MaxUint)),
			}
			// x must be "strict": x ÷ β < y
			x := [3]uint{
				hegel.Draw(ht, hegel.Integers[uint](0, math.MaxUint)),
				hegel.Draw(ht, hegel.Integers[uint](0, math.MaxUint)),
				hegel.Draw(ht, hegel.Integers[uint](0, math.MaxUint)),
			}
			ht.Assume(IsSmaller(x[1:], y[:]))
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			iy := Reciprocal2(y)
			q, r := divisionNormalStrict3By2(x, y, iy)

			// compute with math/big

			x_ := toBigInt(x[:])
			y_ := toBigInt(y[:])
			q_ := big.NewInt(0).Div(x_, y_)
			r_ := big.NewInt(0).Mod(x_, y_)

			// compare

			if big.NewInt(0).SetUint64(uint64(q)).Cmp(q_) != 0 {
				ht.Fatalf("q = %X, q_ = %X", q, q_)
			}
			if toBigInt(r[:]).Cmp(r_) != 0 {
				ht.Fatalf("r = %X, r_ = %X", r, r_)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestDivisionNormalStrict3By2Rapid(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			y := [2]uint(normalLongInteger(rapid.Uint(), 2, 2).Draw(t, "y"))
			t.Logf("y = %X", y)
			isStrict := func(i []uint) bool {
				return IsSmaller(i[1:3], y[:])
			}
			x := [3]uint(rapid.SliceOfN(rapid.Uint(), 3, 3).Filter(isStrict).Draw(t, "x"))
			t.Logf("x = %X", x)

			// compute with purple
			iy := Reciprocal2(y)
			q, r := divisionNormalStrict3By2(x, y, iy)
			t.Logf("q = %X, r = %X", q, r)

			// compute with math/big
			x_ := toBigInt(x[:])
			y_ := toBigInt(y[:])
			q_ := big.NewInt(0).Div(x_, y_)
			r_ := big.NewInt(0).Mod(x_, y_)
			t.Logf("q_ = %X, r_ = %X", q_, r_)

			// compare
			if big.NewInt(0).SetUint64(uint64(q)).Cmp(q_) != 0 {
				t.Error("difference in quotient")
			}
			if toBigInt(r[:]).Cmp(r_) != 0 {
				t.Error("difference in remainder")
			}
		})
	})
}

func TestDivisionNormalStrictN1ByNHegel(t *testing.T) {
	t.Run("accumulate", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			// y must be "normal"
			y := hegel.Draw(ht, hegelNormalLongInteger(1, 63))
			yz := len(y)
			// x must be "strict": x ÷ β < y
			x := hegel.Draw(ht, hegelLongInteger2(yz+1, yz+1))
			ht.Assume(IsSmaller(x[1:], y[:]))
			ht.Logf("x = %X, y = %X", x, y)

			// result

			iy := Reciprocal(y[yz-1])
			r1 := make([]uint, len(x))
			q1 := divisionNormalStrictN1ByN(r1, x, y, iy)

			// accumulate result

			r2 := make([]uint, len(x))
			copy(r2, x)
			q2 := divisionNormalStrictN1ByN(r2, r2, y, iy)

			// compare

			if q1 != q2 {
				ht.Fatalf("q1 = %X, q2 = %X", q1, q2)
			}
			if NotEqual(r1, r2) {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			// y must be "normal"
			y := hegel.Draw(ht, hegelNormalLongInteger(1, 63))
			yz := len(y)
			// x must be "strict": x ÷ β < y
			x := hegel.Draw(ht, hegelLongInteger2(yz+1, yz+1))
			ht.Assume(IsSmaller(x[1:], y[:]))
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			iy := Reciprocal(y[yz-1])
			r := make([]uint, len(x))
			q := divisionNormalStrictN1ByN(r, x, y, iy)

			// compute with math/big

			x_ := toBigInt(x[:])
			y_ := toBigInt(y[:])
			q_ := big.NewInt(0).Div(x_, y_)
			r_ := big.NewInt(0).Mod(x_, y_)

			// compare

			if big.NewInt(0).SetUint64(uint64(q)).Cmp(q_) != 0 {
				ht.Fatalf("q = %X, q_ = %X", q, q_)
			}
			if toBigInt(r).Cmp(r_) != 0 {
				ht.Fatalf("r = %X, r_ = %X", r, r_)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestDivisionNormalStrictN1ByNRapid(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			N := rapid.IntRange(2, 32).Draw(t, "N")
			y := normalLongInteger(rapid.Uint(), N, N).Draw(t, "y")
			isStrict := func(i []uint) bool {
				return IsSmaller(i[1:], y[:])
			}
			x := rapid.SliceOfN(rapid.Uint(), N+1, N+1).Filter(isStrict).Draw(t, "x")

			// compute with purple
			iy := Reciprocal(y[len(y)-1])
			r := make([]uint, len(x))
			q := divisionNormalStrictN1ByN(r, x, y, iy)
			t.Logf("q = %X, r = %X", q, r)

			// compute with math/big
			x_ := toBigInt(x)
			y_ := toBigInt(y)
			q_ := big.NewInt(0).Div(x_, y_)
			r_ := big.NewInt(0).Mod(x_, y_)
			t.Logf("q_ = %X, r_ = %X", q_, r_)

			// compare
			if big.NewInt(0).SetUint64(uint64(q)).Cmp(q_) != 0 {
				t.Error("difference in quotient")
			}
			if toBigInt(r).Cmp(r_) != 0 {
				t.Error("difference in remainder")
			}
		})
	})
	t.Run("accumulate", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			N := rapid.IntRange(2, 32).Draw(t, "N")
			t.Logf("N = %v", N)
			y := normalLongInteger(rapid.Uint(), N, N).Filter(IsNormal).Draw(t, "y")
			t.Logf("y = %v", y)
			isStrict := func(i []uint) bool {
				return IsSmaller(i[1:], y[:])
			}
			x := rapid.SliceOfN(rapid.Uint(), N+1, N+1).Filter(isStrict).Draw(t, "x")
			t.Logf("x = %v", x)

			iy := Reciprocal(y[len(y)-1])

			// compute in result style
			r1 := make([]uint, len(x))
			q1 := divisionNormalStrictN1ByN(r1, x, y, iy)
			t.Logf("q = %v, r = %v", q1, r1)

			// compute in accumulation style
			r2 := make([]uint, len(x))
			copy(r2, x)
			q2 := divisionNormalStrictN1ByN(r2, r2, y, iy)
			t.Logf("q2 = %v, r2 = %v", q2, r2)

			// compare
			if q1 != q2 {
				t.Error("difference in quotient")
			}
			if !slices.Equal(r1, r2) {
				t.Error("difference in remainder")
			}
		})
	})
}

func TestDivisionHegel(t *testing.T) {
	t.Run("accumulate-remainder", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger2(0, 64))
			y := hegel.Draw[[]uint](ht, hegelLongInteger2(1, 64))
			ht.Assume(NotZero(y))
			ht.Logf("x = %X, y = %X", x, y)

			// compute result

			q1 := make([]uint, len(x))
			r1 := make([]uint, len(x))
			Division(q1, r1, x, y)

			// accumulate result

			q2 := make([]uint, len(x))
			r2 := make([]uint, len(x))
			copy(r2, x)
			Division(q2, r2, r2, y)

			// compare

			if NotEqual(q1, q2) {
				ht.Fatalf("q1 = %X, q2 = %X", q1, q2)
			}
			if NotEqual(r1, r2) {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("differential", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger2(0, 64))
			y := hegel.Draw[[]uint](ht, hegelLongInteger2(1, 64))
			ht.Assume(NotZero(y))
			ht.Logf("x = %X, y = %X", x, y)

			// compute with purple

			q := make([]uint, len(x))
			r := make([]uint, len(y))
			Division(q, r, x, y)

			// compute with math/big

			x_ := toBigInt(x)
			y_ := toBigInt(y)
			q_ := big.NewInt(0).Div(x_, y_)
			r_ := big.NewInt(0).Mod(x_, y_)

			// compare

			if toBigInt(q).Cmp(q_) != 0 {
				ht.Fatalf("q = %X, q_ = %X", q, q_)
			}
			if toBigInt(r).Cmp(r_) != 0 {
				ht.Fatalf("r = %X, r_ = %X", r, r_)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("short-quotient", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger2(1, 64))
			y := hegel.Draw[[]uint](ht, hegelLongInteger2(1, 64))
			ht.Assume(NotZero(y))
			// full size
			fz := len(x)
			// short size
			sz := hegel.Draw(ht, hegel.Integers[int](0, fz-1))
			ht.Logf("x = %X, y = %X, sz = %d", x, y, sz)

			// full quotient

			q1 := make([]uint, fz)
			r1 := make([]uint, len(y))
			Division(q1, r1, x, y)

			// short quotient

			q2 := make([]uint, sz)
			r2 := make([]uint, len(y))
			Division(q2, r2, x, y)

			// compare

			if NotEqual(q1[:sz], q2) {
				ht.Fatalf("q1 = %X, q2 = %X", q1, q2)
			}
			if NotEqual(r1, r2) {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
	t.Run("short-remainder", func(t *testing.T) {
		hegel.Test(t, func(ht *hegel.T) {

			// generate samples

			x := hegel.Draw[[]uint](ht, hegelLongInteger2(1, 64))
			y := hegel.Draw[[]uint](ht, hegelLongInteger2(1, 64))
			ht.Assume(NotZero(y))
			// full size
			fz := len(x)
			// short size
			sz := hegel.Draw(ht, hegel.Integers[int](0, fz-1))
			ht.Logf("x = %X, y = %X, sz = %d", x, y, sz)

			// full remainder

			q1 := make([]uint, len(x))
			r1 := make([]uint, fz)
			Division(q1, r1, x, y)

			// short remainder

			q2 := make([]uint, len(x))
			r2 := make([]uint, sz)
			Division(q2, r2, x, y)

			// compare

			if NotEqual(q1, q2) {
				ht.Fatalf("q1 = %X, q2 = %X", q1, q2)
			}
			if NotEqual(r1[:sz], r2) {
				ht.Fatalf("r1 = %X, r2 = %X", r1, r2)
			}

		}, hegel.WithTestCases(hegelCases))
	})
}

func TestDivisionRapid(t *testing.T) {
	t.Run("differential", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// generate samples
			y := rapid.SliceOfN(rapid.Uint(), 2, 32).Filter(IsCompact).Filter(NotZero).Draw(t, "y")
			x := rapid.SliceOfN(rapid.Uint(), 0, 32).Draw(t, "x")

			// compute with purple
			q := make([]uint, len(x))
			r := make([]uint, len(y))
			Division(q, r, x, y)
			t.Logf("q = %X, r = %X", q, r)

			// compute with math/big
			x_ := toBigInt(x)
			y_ := toBigInt(y)
			t.Logf("x_ = %X, y_ = %X", x_, y_)
			q_ := big.NewInt(0).Div(x_, y_)
			r_ := big.NewInt(0).Mod(x_, y_)
			t.Logf("q_ = %X, r_ = %X", q_, r_)

			// compare
			if toBigInt(q).Cmp(q_) != 0 {
				t.Error("difference in quotient")
			}
			if toBigInt(r).Cmp(r_) != 0 {
				t.Error("difference in remainder")
			}
		})
	})
}

func BenchmarkDivisionNormalStrict2By1(b *testing.B) {
	const Bits = bits.UintSize
	rng := newRand()

	// generate samples
	y := rng.Uint()
	y |= 1 << (Bits - 1)
	var x [2]uint
	x[0] = rng.Uint()
	x[1] = rng.Uint()
	if x[1] >= y {
		x[1] = x[1] - y
	}

	// measure purple
	b.Run("purple", func(b *testing.B) {
		iy := Reciprocal(y)
		var q, r uint
		for b.Loop() {
			q, r = divisionNormalStrict2By1(x, y, iy)
		}
		_, _ = q, r
	})

	// measure bits
	b.Run("bits", func(b *testing.B) {
		var q, r uint
		for b.Loop() {
			q, r = bits.Div(x[1], x[0], y)
		}
		_, _ = q, r
	})
}

func BenchmarkDivisionNormalStrict3By2(b *testing.B) {
	const Bits = bits.UintSize
	rng := newRand()

	// generate samples
	var y [2]uint
	y[0] = rng.Uint()
	y[1] = rng.Uint()
	y[1] |= 1 << (Bits - 1)
	var x [3]uint
	x[0] = rng.Uint()
	x[1] = rng.Uint()
	x[2] = rng.Uint()
	if NotSmaller(x[1:], y[:]) {
		_ = Difference(x[1:], x[1:], y[:])
	}

	// measure purple
	b.Run("purple", func(b *testing.B) {
		iy := Reciprocal2(y)
		var q uint
		var r [2]uint
		for b.Loop() {
			q, r = divisionNormalStrict3By2(x, y, iy)
		}
		_, _ = q, r
	})

	// measure math/big
	b.Run("math/big", func(b *testing.B) {
		x_ := toBigInt(x[:])
		y_ := toBigInt(y[:])
		q := big.NewInt(0)
		r := big.NewInt(0)
		for b.Loop() {
			q, r = q.DivMod(x_, y_, r)
		}
		_, _ = q, r
	})
}

func BenchmarkDivisionNormalStrictN1ByN(b *testing.B) {
	const Bits = bits.UintSize
	rng := newRand()

	for _, N := range []uint{8, 16, 32, 64, 128, 256} {
		// generate samples
		y := make([]uint, N)
		for i := range y {
			y[i] = rng.Uint()
		}
		y[N-1] |= 1 << (Bits - 1)
		x := make([]uint, N+1)
		for i := range x {
			x[i] = rng.Uint()
		}
		if NotSmaller(x[1:], y[:]) {
			_ = Difference(x[1:], x[1:], y[:])
		}

		// measure purple
		b.Run(fmt.Sprint("purple-", N), func(b *testing.B) {
			iy := Reciprocal(y[N-1])
			var q uint
			r := make([]uint, N)
			for b.Loop() {
				q = divisionNormalStrictN1ByN(r[:], x[:], y[:], iy)
			}
			_, _ = q, r
		})

		// measure math/big
		b.Run(fmt.Sprint("math-big-", N), func(b *testing.B) {
			x_ := toBigInt(x[:])
			y_ := toBigInt(y[:])
			q := big.NewInt(0)
			r := big.NewInt(0)
			for b.Loop() {
				q, r = q.DivMod(x_, y_, r)
			}
			_, _ = q, r
		})
	}
}

func BenchmarkDivision(b *testing.B) {
	const Bits = bits.UintSize
	rng := newRand()

	for _, N := range []uint{8, 16, 32, 64, 128, 256} {
		// generate samples
		y := make([]uint, N)
		for i := range y {
			y[i] = rng.Uint()
		}
		for y[N-1] == 0 {
			y[N-1] = rng.Uint()
		}
		x := make([]uint, N)
		for i := range x {
			x[i] = rng.Uint()
		}

		// measure purple
		b.Run(fmt.Sprint("purple-", N), func(b *testing.B) {
			q := make([]uint, N)
			r := make([]uint, N)
			for b.Loop() {
				Division(q[:], r[:], x[:], y[:])
			}
			_, _ = q, r
		})

		// measure math/big
		b.Run(fmt.Sprint("math-big-", N), func(b *testing.B) {
			x_ := toBigInt(x[:])
			y_ := toBigInt(y[:])
			q := big.NewInt(0)
			r := big.NewInt(0)
			for b.Loop() {
				q, r = q.DivMod(x_, y_, r)
			}
			_, _ = q, r
		})
	}
}
