// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// divisionNormalStrict2By1 computes the ratio of a two-word integer by a one-word integer.
//
// Requires:
// y is normalized,
// iy = Reciprocal(y),
// x[1] < y;
// otherwise, the result will be wrong.
//
// This implementation applies the "Improved division by invariant integers" method.
func divisionNormalStrict2By1(x [2]uint, y uint, iy uint) (quotient uint, remainder uint) {
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

// DivisionBy1 computes the ratio of a multi-word integer by a one-word integer.
//
// DivisionBy1 requires:
// y is nonzero;
// otherwise, it gives the wrong result.
//
// DivisionBy1 adds into quotient the len(quotient) words of the result.
// It permits aliasing quotient to x, in which case it becomes "divide and add".
//
// This implementation applies the "Improved division by invariant integers" method.
func DivisionBy1(quotient []uint, x []uint, y uint) (remainder uint) {
	qz := len(quotient)
	xz := len(x)

	if !(qz >= xz) {
		panic("requires len(quotient) >= len(x)")
	}

	if xz == 0 {
		return
	}

	// dividend buffer
	tz := xz + 1
	t := make([]uint, tz)
	copy(t, x)

	// normalize operands
	factor := uint(bits.LeadingZeros(y))
	y = y << factor
	_ = Twice(t, t, factor)
	// invariant: t did not overflow

	// compute reciprocal approximation
	iy := Reciprocal(y)

	// compute quotient and remainder, word by word
	for i := tz - 1; i > 0; i-- {
		x_ := [2]uint(t[i-1 : i+1])
		// invariant: x_[i] < y
		q, r := divisionNormalStrict2By1(x_, y, iy)
		quotient[i-1] += q
		t[i-1] = r
	}

	// de-"normalise" remainder
	remainder = t[0] >> factor
	// invariant: remainder did not underflow

	return
}

// divisionNormalStrict3By2 computes the ratio of a three-word integer by a two-word integer.
//
// divisionNormalStrict3By2 requires:
// y is normalized,
// iy = Reciprocal2(y[1:]),
// x[1:3] < y;
// otherwise, the result will be wrong.
//
// This implementation applies the "Improved division by invariant integers" method.
func divisionNormalStrict3By2(x [3]uint, y [2]uint, iy uint) (quotient uint, remainder [2]uint) {
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

// divisionNormalStrictN1ByN computes the ratio of an (N+1)-word integer by a N-word integer.
//
// divisionNormalStrictN1ByN requires:
// len(y) > 1,
// len(x) = len(y) + 1,
// y is normalized,
// iy is Reciprocal(y[len(y)-1]),
// x[1:] < y;
// otherwise, the result will be wrong.
//
// divisionNormalStrictN1ByN permits aliasing remainder to x.
//
// This implementation applies the "school" method described in Knuth, section 4.3.1, steps D3 through D6,
// enhanced by the "Improved division by invariant integers" at step D3.
func divisionNormalStrictN1ByN(remainder []uint, x []uint, y []uint, iy uint) (quotient uint) {
	yz := len(y)

	// step D3: calculate q'
	// q' ← ( u[j+n]×β + u[j+n-1] ) ÷ y[n-1]
	// r' ← ( u[j+n]×β + u[j+n-1] ) % y[n-1]
	var q_ [2]uint
	var r_ uint
	q_[1], r_ = divisionNormalStrict2By1([2]uint{x[yz], r_}, y[yz-1], iy)
	q_[0], r_ = divisionNormalStrict2By1([2]uint{x[yz-1], r_}, y[yz-1], iy)
	// invariant: q' - 2 ≤ quotient ≤ q' ≤ β

	// step D3: reduce q'
	test := func(q_ [2]uint, x []uint, y []uint, r_ uint) bool {
		// let t0 = q' × v[n-2]
		var t0 [2]uint
		t0[1], t0[0] = bits.Mul(q_[0], y[yz-2])
		// let t1 = r'×β + u[j+n-2]
		var t1 [2]uint
		t1 = [2]uint{x[yz-2], r_}
		// test q' × v[n-2] > r'×β + u[j+n-2]
		return IsGreater(t0[:], t1[:])
	}
	// if q' >= β or q' × v[n-2] > r'×β + u[j+n-2]
	for q_[1] != 0 || test(q_, x, y, r_) {
		// then fix q', r'
		_ = Difference(q_[:], q_[:], []uint{1})
		var carry uint
		r_, carry = bits.Add(r_, y[yz-1], 0)
		// if r' < β, repeat
		if carry != 0 {
			break
		}
	}
	// invariant: q' - 1 ≤ q ≤ q' < β

	// step D4: multiply and subtract
	var borrow uint
	{
		// let t = q' × y
		t := make([]uint, yz+1)
		t[yz] = ProductBy1(t, y, q_[0])
		// remainder ← x - q' × y
		borrow = Difference(remainder, x, t)
	}

	// step D5: test remainder
	if borrow != 0 {
		// step D6: add back
		q_[0], _ = bits.Sub(q_[0], 1, 0)
		_ = Sum(remainder, remainder, y)
	}
	// invariant: q' = q

	return q_[0]
}

// Division computes the ratio of an M-word integer by an N-word integer.
//
// Requires:
// y is compact;
// y is not zero;
// otherwise, it gives the wrong result.
//
// This implementation applies the "school" method described in Knuth, section 4.3.1,
// enhanced by the "Improved division by invariant integers".
func Division(quotient []uint, remainder []uint, x []uint, y []uint) {
	qz := len(quotient)
	rz := len(remainder)
	xz := len(x)
	yz := len(y)

	if !(yz > 0) {
		panic("requires len(y) > 0")
	}
	if !(qz >= xz) {
		panic("requires len(quotient) >= len(x)")
	}
	if !(rz >= yz) {
		panic("requires len(remainder) >= len(y)")
	}

	if yz == 1 {
		remainder[0] = DivisionBy1(quotient, x, y[0])
		return
	}

	if xz < yz {
		copy(remainder, x)
		return
	}

	// dividend buffer
	tz := xz + 1
	t := make([]uint, tz)
	copy(t, x)

	// normalize operands
	factor := uint(bits.LeadingZeros(y[yz-1]))
	_ = Twice(y, y, factor)
	_ = Twice(t, t, factor)
	// invariant: t did not overflow

	// compute reciprocal approximation of topmost dividend word
	iy := Reciprocal(y[yz-1])

	// compute quotient and remainder, word by word
	for i := tz - yz; i > 0; i-- {
		x_ := t[i-1 : i+yz]
		// invariant: x_[1:] < y
		quotient[i-1] += divisionNormalStrictN1ByN(x_, x_, y, iy)
	}

	// denormalize operands
	_ = Half(y, y, factor)
	_ = Half(remainder, t, factor)
	// invariant: remainder did not underflow

	return
}
