package arithmetics

func AssignUni(r []uint, x uint) {
	rz := len(r)

	if 0 < rz {
		r[0] = x
	}

	for i := 1; i < rz; i++ {
		r[i] = 0
	}
}

func Assign(r, x []uint) {
	rz := len(r)
	xz := len(x)

	z := min(rz, xz)

	for i := 0; i < z; i++ {
		r[i] = x[i]
	}

	for i := z; i < rz; i++ {
		r[i] = 0
	}
}
