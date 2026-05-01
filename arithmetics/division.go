// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// divisionNormalStrict2By1 of nonnegative 2-word integer `x` by 1-word normalized integer `y`.
//
// Returns the quotient and the remainder.
//
// Requires:
// y is normalized;
// iy = reciprocal_normalized(y);
// x ÷ β < y.
//
// This implementation applies the "improved division by invariant integers" method.
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

// DivisionUni of nonnegative integers `x` by `y`.
//
// Stores into `q` the `size(q)` least significant words of the quotient.
// Returns the remainder.
//
// Requires:
// y is nonzero.
//
// This implementation applies the "school" method described in Knuth, section 4.3.1,
// augmented by the "improved division by invariant integers" method.
func DivisionUni(quotient []uint, x []uint, y uint) (remainder uint) {
	qz := len(quotient)
	xz := len(x)

	// 1. normalize.

	// 1.1. normalization factor.
	factor := uint(bits.LeadingZeros(y))

	// 1.2. normalize divisor.
	ny := y << factor

	// 1.3. fix dividend, which becomes "strict".
	nx := make([]uint, xz+1)
	nx[xz] = Twice(nx, x, factor)
	// invariant: t[-z] < y

	// 2. reciprocal approximation of normalized divisor.

	iy := Reciprocal(ny)

	// 3. compute quotient, word by word.

	r_ := nx[xz]

	for i := xz; i > min(qz, xz); i-- {
		_, r_ = divisionNormalStrict2By1([2]uint{nx[i-1], r_}, ny, iy)
	}

	for i := min(qz, xz); i > 0; i-- {
		quotient[i-1], r_ = divisionNormalStrict2By1([2]uint{nx[i-1], r_}, ny, iy)
	}

	// 4. denormalize remainder.

	remainder = r_ >> factor
	return
}

// divisionNormalStrict3By2 of nonnegative 3-word integer `x` by 2-word normalized integer `y`.
//
// Returns the quotient and the remainder.
//
// Requires:
// y is normalized
// iy = reciprocal_normalized(y)
// x ÷ β < y
//
// This implementation applies the "improved division by invariant integers" method.
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

// divisionNormalStrictN1ByN of nonnegative (N+1)-word integer `x` by N-word normalized integer `y`.
//
// Stores into `r` the `size(r)` least significant words of the remainder.
// Permits aliasing `r` to `x`, in which case it "accumulates" the remainder.
// Returns the quotient.
//
// Requires:
// size(r) ≥ size(x)
// size(x) = size(y) + 1
// size(y) ≥ 2
// y is normalized
// iy = reciprocal_normalized( top(y) )
// x ÷ β < y
//
// This implementation applies the "school" method described in Knuth, section 4.3.1,
// augmented by the "improved division by invariant integers" method.
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
		t[yz] = ProductUni(t, y, q_[0])
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

// Division of nonnegative integers `x` by `y`.
//
// Stores into `q` the `size(q)` least significant words of the quotient.
// Stores into `r` the `size(r)` least significant words of the remainder.
//
// Requires:
// y is nonzero
//
// This implementation applies the "school" method described in Knuth, section 4.3.1,
// augmented by the "improved division by invariant integers" method.
func Division(q []uint, r []uint, x []uint, y []uint) {
	y = Compact(y)

	qz := len(q)
	//rz := len(remainder)
	xz := len(x)
	yz := len(y)

	if !(yz > 0) {
		panic("requires len(y) > 0")
	}

	if yz == 1 {
		r_ := DivisionUni(q, x, y[0])
		AssignUni(r, r_)
		return
	}
	// invariant: size(y) > 1

	if xz < yz {
		AssignUni(q, 0)
		Assign(r, x)
		return
	}
	// invariant: size(x) >= size(y)

	// 1. normalize

	// 1.1. normalization factor.
	factor := uint(bits.LeadingZeros(y[yz-1]))

	// 1.2. normalize divisor.
	ny := make([]uint, yz)
	_ = Twice(ny, y, factor)

	// 1.3. fix dividend, which becomes "strict".
	nx := make([]uint, xz+1)
	nx[xz] = Twice(nx, x, factor)
	// invariant: size(nx) > size(ny)
	// invariant: nx ÷ β < ny

	// 2. compute reciprocal approximation of normalized divisor top word.

	iy := Reciprocal(ny[yz-1])

	// 3. compute quotient, word by word.

	M := xz + 1 - yz
	// invariant: M ≥ 0

	for i := M; i > min(qz, M); i-- {
		x_ := nx[i-1 : i+yz]
		_ = divisionNormalStrictN1ByN(x_, x_, ny, iy)
	}

	for i := min(qz, M); i > 0; i-- {
		x_ := nx[i-1 : i+yz]
		q[i-1] += divisionNormalStrictN1ByN(x_, x_, ny, iy)
	}

	// 4. denormalize remainder.

	_ = Half(r, nx, factor)

	return
}
