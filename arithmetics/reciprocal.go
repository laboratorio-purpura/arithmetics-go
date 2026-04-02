package arithmetics

import (
	"math"
	"math/bits"
)

// Reciprocal computes an approximation to the multiplicative inverse of a "normalised" one-word integer.
//
// Requires:
// y is "normalised".
// Otherwise, the result is undefined.
//
// This implementation applies the "Improved division by invariant integers" method.
func Reciprocal(y uint) uint {
	// t ← β - 1
	t := uint(math.MaxUint)
	// reciprocal ← <β - 1 - y, β - 1> ÷ y
	iy, _ := bits.Div(t-y, t, y)
	return iy
}

// Reciprocal2 computes an approximation to the multiplicative inverse of a "normalised" two-word integer.
//
// Requires:
// y is "normalised".
// Otherwise, the result is undefined.
//
// This implementation applies the "Improved division by invariant integers" method.
func Reciprocal2(y [2]uint) uint {
	// 1. v ← RECIPROCAL_WORD(d1)
	v := Reciprocal(y[1])
	// We have β^2 − d1 ≤ (β + v).d1 < β^2
	// 2. p ← d1.v mod β
	_, p := bits.Mul(y[1], v)
	// 3. p ← (p + d0) mod β
	p, _ = bits.Add(p, y[0], 0)
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
		p, _ = bits.Sub(p, y[1], 0)
	}
	// We have β^2 − d1 ≤ (β + v) . d1 + d0 < β^2.
	// 10. <t1, t0> ← v.d0
	t1, t0 := bits.Mul(v, y[0])
	// 11. p ← (p + t1) mod β
	p, _ = bits.Add(p, t1, 0)
	// 12. if p < t1
	if p < t1 {
		// 13. v ← v − 1
		v = v - 1
		// 14. if <p, t0> ≥ <d1, d0>
		if NotSmaller([]uint{t0, p}, y[:]) {
			// 15. v ← v − 1
			v = v - 1
		}
	}
	return v
}
