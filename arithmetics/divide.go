// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// DivideNormalStrict2By1 computes the ratio of a two-word integer by a one-word integer.
//
// Requires:
// y is normalized,
// iy = Reciprocal(y),
// x[1] < y;
// otherwise, the result will be wrong.
//
// This implementation applies the "Improved division by invariant integers" method.
func DivideNormalStrict2By1(x [2]uint, y uint, iy uint) (quotient uint, remainder uint) {
	// 1. <q1, q0> ← v.u1
	q1, q0 := bits.Mul(x[1], iy)
	// 2. <q1, q0> ← <q1, q0> + <u1, u0>
	q0, carry := bits.Add(q0, x[0], 0)
	q1, _ = bits.Add(q1, x[1], carry)
	// 3. q1 ← (q1 + 1) mod β
	q1, _ = bits.Add(q1, 1, 0)
	// 4. r ← (u0 − q1.d) mod β
	_, t := bits.Mul(q1, y)
	r, _ := bits.Sub(x[0], t, 0)
	// 5. if r > q0
	if r > q0 {
		// 6. q1 ← (q1 − 1) mod β
		q1, _ = bits.Sub(q1, 1, 0)
		// r ← (r + d) mod β
		r, _ = bits.Add(r, y, 0)
	}
	// 8. if r ≥ d
	if r >= y {
		// 9. q1 ← q1 + 1
		q1, _ = bits.Add(q1, 1, 0)
		// 10. r ← r − d
		r, _ = bits.Sub(r, y, 0)
	}
	return q1, r
}

// DivideNormalNBy1 computes the ratio of an N-word integer by a one-word integer.
//
// DivideNormalNBy1 requires:
// y is normalized,
// iy = Reciprocal(y);
// otherwise, the result will be wrong.
//
// DivideNBy1 adds into quotient the len(quotient) words of the result.
// It permits aliasing quotient to x, in which case it becomes "divide and add".
//
// This implementation applies the "Improved division by invariant integers" method.
func DivideNormalNBy1(quotient []uint, x []uint, y uint, iy uint) (remainder uint) {
	xz := len(x)

	if len(quotient) < xz {
		panic("")
	}

	// dividend buffer
	tz := xz + 1
	t := make([]uint, tz)
	copy(t, x)
	// invariant: t[tz-1] < y

	// compute quotient and remainder, word by word
	for i := tz - 1; i > 0; i-- {
		x_ := [2]uint(t[i-1 : i+1])
		// invariant: x_[i] < y
		q, r := DivideNormalStrict2By1(x_, y, iy)
		quotient[i-1] += q
		t[i-1] = r
	}

	return t[0]
}

// DivideNBy1 computes the ratio of a multi-word integer by a one-word integer.
//
// DivideNBy1 adds into quotient the len(quotient) words of the result.
// It permits aliasing quotient to x, in which case it becomes "divide and add".
//
// This implementation applies the "Improved division by invariant integers" method.
func DivideNBy1(quotient []uint, x []uint, y uint) (remainder uint) {
	qz := len(quotient)
	xz := len(x)

	if qz < xz {
		panic("requires len(quotient) >= len(x)")
	}

	if xz == 0 {
		return
	}
	// invariant: xz >= 1

	// dividend buffer
	tz := xz + 1
	t := make([]uint, tz)
	copy(t, x)
	// invariant: tz > z

	// "normalise" operands
	factor := uint(bits.LeadingZeros(y))
	y = y << factor
	_ = Double(t, t, factor)
	// invariant: y is "normalised"
	// invariant: t did not overflow

	// compute reciprocal approximation
	iy := Reciprocal(y)

	// compute quotient and remainder, word by word
	for i := tz - 1; i > 0; i-- {
		x_ := [2]uint(t[i-1 : i+1])
		// invariant: x_[i] < y
		q, r := DivideNormalStrict2By1(x_, y, iy)
		quotient[i-1] += q
		t[i-1] = r
	}

	// de-"normalise" remainder
	remainder = t[0] >> factor
	// invariant: remainder did not underflow

	return
}

// DivideNormal3By2WithReciprocal computes the ratio of a three-word integer by a two-word integer.
//
// Requires:
// y is "normalised";
// iy == Reciprocal2(y);
// x[2:3] < y;
// otherwise, the result will be wrong.
//
// This implementation applies the "Improved division by invariant integers" method.
func DivideNormal3By2WithReciprocal(x [3]uint, y [2]uint, iy uint) (quotient uint, remainder [2]uint) {
	// 1. <q1,q0> ← v.u2
	q1, q0 := bits.Mul(iy, x[2])
	// 2. <q1,q0> ← <q1,q0> + <u2,u1>
	q0, carry := bits.Add(q0, x[1], 0)
	q1, _ = bits.Add(q1, x[2], carry)
	// 3. r1 ← (u1 − q1.d1) mod β
	_, t := bits.Mul(q1, y[1])
	r1, _ := bits.Sub(x[1], t, 0)
	// 4. <t1,t0> ← d0.q1
	t1, t0 := bits.Mul(y[0], q1)
	// 5. <r1,r0> ← (<r1,u0> − <t1,t0> − <d1,d0>) mod β^2
	r0, borrow := bits.Sub(x[0], t0, 0)
	r1, _ = bits.Sub(r1, t1, borrow)
	r0, borrow = bits.Sub(r0, y[0], 0)
	r1, _ = bits.Sub(r1, y[1], borrow)
	// q1 ← (q1 + 1) mod β
	q1, _ = bits.Add(q1, 1, 0)
	// 7. if r1 ≥ q0
	if r1 >= q0 {
		// 8. q1 ← (q1 − 1) mod β
		q1, _ = bits.Sub(q1, 1, 0)
		// 9. <r1,r0> ← (<r1,r0> + <d1,d0>) mod β^2
		r0, carry = bits.Add(r0, y[0], 0)
		r1, _ = bits.Add(r1, y[1], carry)
	}
	// 10. if <r1,r0> ≥ <d1,d0>
	if NotSmaller([]uint{r0, r1}, y[:]) {
		// 11. q1 ← q1 + 1
		q1, _ = bits.Add(q1, 1, 0)
		// 12. <r1,r0> ← <r1,r0> − <d1,d0>
		r0, borrow = bits.Sub(r0, y[0], 0)
		r1, _ = bits.Sub(r1, y[1], borrow)
	}
	return q1, [2]uint{r0, r1}
}

// DivideNormalN1ByN computes the ratio of an (N+1)-word integer by a N-word integer.
//
// DivideNormalN1ByN requires:
// len(y) > 1;
// len(x) = len(y) + 1;
// y is "normalised";
// x[N-1:N+1] < y[N-1];
// otherwise, the result will be wrong.
//
// DivideNormalN1ByN permits aliasing remainder to x, in which case it becomes "divide accumulate".
//
// This implementation applies the "school" method described in Knuth, section 4.3.1, steps D3 through D6.
// optimized by eliding a test in step D3.
func DivideNormalN1ByN(remainder []uint, x []uint, y []uint) (quotient uint) {
	xz := len(x)
	yz := len(y)

	// step D3: calculate q'
	q_, r_ := bits.Div(x[xz-1], x[xz-2], y[yz-1])
	// invariant: q' - 2 ≤ quotient ≤ q'

	// step D3: reduce q'
	for {
		// if q' >= β…
		// but it will never,
		// because we require x[N-1:N+1] < y[N-1] => x[N-1:N+1] / y[N-1] < β.

		// or if q' × y[yz-2] > { x[yz-2], r' }…
		// let t0 = q' × y[yz-2]
		var t0 [2]uint
		t0[1], t0[0] = bits.Mul(q_, y[yz-2])
		// let t1 = { x[yz-2], r' }
		var t1 [2]uint
		t1 = [2]uint{x[xz-2], r_}
		if NotGreater(t0[:], t1[:]) {
			break
		}

		// then fix q', r'
		q_ = q_ - 1
		var carry uint
		r_, carry = bits.Add(r_, y[yz-1], 0)

		// if r' < β, repeat
		if carry != 0 {
			break
		}
	}
	// invariant: q' - 1 ≤ q ≤ q'

	// step D4: multiply and subtract
	var borrow uint
	{
		// let t = q' × y
		t := make([]uint, yz+1)
		Multiply(t, []uint{q_}, y)
		// remainder <- x - q' × y
		borrow = Subtract(remainder, x, t)
	}

	// step D5: test remainder
	if borrow != 0 {
		// step D6: add back
		q_ = q_ - 1
		_ = Add(remainder, remainder, y)
	}
	// invariant: q' = q

	return q_
}
