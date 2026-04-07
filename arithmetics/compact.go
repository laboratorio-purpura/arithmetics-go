// Copyright (c) 2026 Pedro Lamarão <pedro.lamarao@purpura.dev.br>
// SPDX-License-Identifier: GPL-3.0-only

package arithmetics

func IsCompact(x []uint) bool {
	return len(x) < 2 || x[len(x)-1] != 0
}

func NotCompact(x []uint) bool {
	return len(x) >= 2 && x[len(x)-1] == 0
}

func Compact(x []uint) []uint {
	for NotCompact(x) {
		x = x[:len(x)-1]
	}
	return x
}
