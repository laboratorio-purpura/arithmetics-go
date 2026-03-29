package arithmetics

import (
	"math"
	"math/bits"
)

func Reciprocal1W32(y uint32) uint32 {
	// t ← β - 1
	t := uint32(math.MaxUint32)
	// reciprocal ← <β - 1 - y, β - 1> ÷ y
	iy, _ := bits.Div32(t-y, t, y)
	return iy
}

func Reciprocal2W32(y [2]uint32) uint32 {
	// 1. v ← RECIPROCAL_WORD(d1)
	v := Reciprocal1W32(y[1])
	// We have β^2 − d1 ≤ (β + v).d1 < β^2
	// 2. p ← d1.v mod β
	_, p := bits.Mul32(y[1], v)
	// 3. p ← (p + d0) mod β
	p, _ = bits.Add32(p, y[0], 0)
	// 4. if p < d0
	if p < y[0] {
		// 5. v ← v − 1
		v = v - 1
		// 6. if p ≥ d1
		if p >= y[1] {
			// 7. v ← v − 1
			v = v - 1
			// 8. p ← p − d1
			p = p - y[1]
		}
		// 9. p ← (p − d1) mod β
		p, _ = bits.Sub32(p, y[1], 0)
	}
	// We have β^2 − d1 ≤ (β + v) . d1 + d0 < β^2.
	// 10. <t1, t0> ← v.d0
	t1, t0 := bits.Mul32(v, y[0])
	// 11. p ← (p + t1) mod β
	p, _ = bits.Add32(p, t1, 0)
	// 12. if p < t1
	if p < t1 {
		// 13. v ← v − 1
		v = v - 1
		// 14. if <p, t0> ≥ <d1, d0>
		if NotLess32([]uint32{t0, p}, y[:]) {
			// 15. v ← v − 1
			v = v - 1
		}
	}
	return v
}
