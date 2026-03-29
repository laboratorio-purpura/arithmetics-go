// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import (
	"math/bits"
)

// Division2By1WithReciprocal32 computes the division of a two-word integer by a one-word integer.
//
// Computes by the "Improved division by invariant integers" method.
func Division2By1WithReciprocal32(x [2]uint32, y uint32, iy uint32) (quotient uint32, remainder uint32) {
	// 1. <q1, q0> ← v.u1
	q1, q0 := bits.Mul32(x[1], iy)
	// 2. <q1, q0> ← <q1, q0> + <u1, u0>
	q0, carry := bits.Add32(q0, x[0], 0)
	q1, _ = bits.Add32(q1, x[1], carry)
	// 3. q1 ← (q1 + 1) mod β
	q1, _ = bits.Add32(q1, 1, 0)
	// 4. r ← (u0 − q1.d) mod β
	_, t := bits.Mul32(q1, y)
	r, _ := bits.Sub32(x[0], t, 0)
	// 5. if r > q0
	if r > q0 {
		// 6. q1 ← (q1 − 1) mod β
		q1, _ = bits.Sub32(q1, 1, 0)
		// r ← (r + d) mod β
		r, _ = bits.Add32(r, y, 0)
	}
	// 8. if r ≥ d
	if r >= y {
		// 9. q1 ← q1 + 1
		q1, _ = bits.Add32(q1, 1, 0)
		// 10. r ← r − d
		r, _ = bits.Sub32(r, y, 0)
	}
	return q1, r
}

// DivisionBy1WithReciprocal32 computes the division of a "long" integer by a one-word integer.
//
// Computes by the "Improved division by invariant integers" method.
func DivisionBy1WithReciprocal32(q []uint32, x []uint32, y uint32, iy uint32) (r uint32) {
	xz := len(x)
	for i := xz - 1; i >= 0; i-- {
		x_ := [2]uint32{x[i], r}
		q[i], r = Division2By1WithReciprocal32(x_, y, iy)
	}
	return r
}

// Division3By2WithReciprocal32 computes the division of a three-word integer by a two-word integer.
//
// Computes by the "Improved division by invariant integers" method.
func Division3By2WithReciprocal32(x [3]uint32, y [2]uint32, iy uint32) (q uint32, r [2]uint32) {
	// 1. <q1,q0> ← v.u2
	q1, q0 := bits.Mul32(iy, x[2])
	// 2. <q1,q0> ← <q1,q0> + <u2,u1>
	q0, carry := bits.Add32(q0, x[1], 0)
	q1, _ = bits.Add32(q1, x[2], carry)
	// 3. r1 ← (u1 − q1.d1) mod β
	_, t := bits.Mul32(q1, y[1])
	r1, _ := bits.Sub32(x[1], t, 0)
	// 4. <t1,t0> ← d0.q1
	t1, t0 := bits.Mul32(y[0], q1)
	// 5. <r1,r0> ← (<r1,u0> − <t1,t0> − <d1,d0>) mod β^2
	r0, borrow := bits.Sub32(x[0], t0, 0)
	r1, _ = bits.Sub32(r1, t1, borrow)
	r0, borrow = bits.Sub32(r0, y[0], 0)
	r1, _ = bits.Sub32(r1, y[1], borrow)
	// q1 ← (q1 + 1) mod β
	q1, _ = bits.Add32(q1, 1, 0)
	// 7. if r1 ≥ q0
	if r1 >= q0 {
		// 8. q1 ← (q1 − 1) mod β
		q1, _ = bits.Sub32(q1, 1, 0)
		// 9. <r1,r0> ← (<r1,r0> + <d1,d0>) mod β^2
		r0, carry = bits.Add32(r0, y[0], 0)
		r1, _ = bits.Add32(r1, y[1], carry)
	}
	// 10. if <r1,r0> ≥ <d1,d0>
	if NotLess32([]uint32{r0, r1}, y[:]) {
		// 11. q1 ← q1 + 1
		q1, _ = bits.Add32(q1, 1, 0)
		// 12. <r1,r0> ← <r1,r0> − <d1,d0>
		r0, borrow = bits.Sub32(r0, y[0], 0)
		r1, _ = bits.Sub32(r1, y[1], borrow)
	}
	return q1, [2]uint32{r0, r1}
}
