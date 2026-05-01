// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

func IsCompact(x []uint) bool {
	return len(x) < 2 || x[len(x)-1] != 0
}

func NotCompact(x []uint) bool {
	return len(x) == 1 || x[len(x)-1] == 0
}

func Compact(x []uint) []uint {
	i := len(x)
	for i > 1 && x[i-1] == 0 {
		i--
	}
	return x[:i]
}
