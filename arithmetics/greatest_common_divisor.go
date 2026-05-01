package arithmetics

func GreatestCommonDivisorClobber(x, y []uint) []uint {
	for NotZero(y) {
		Division([]uint{}, x, x, y)
		x, y = y, x
	}
	return x
}
