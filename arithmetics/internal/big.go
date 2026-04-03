package internal

import (
	"math/big"
	"math/bits"
)

func ToBigInt(x []uint) *big.Int {
	x_ := big.NewInt(0)
	for i := len(x); i > 0; i-- {
		w := big.NewInt(0).SetUint64(uint64(x[i-1]))
		x_.Lsh(x_, bits.UintSize).Add(x_, w)
	}
	return x_
}
