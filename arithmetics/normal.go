// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

import "math/bits"

func IsNormal(x []uint) bool {
	return len(x) > 0 && x[len(x)-1]&(1<<(bits.UintSize-1)) != 0
}

func NotNormal(x []uint) bool {
	return len(x) == 0 || x[len(x)-1]&(1<<(bits.UintSize-1)) == 0
}
