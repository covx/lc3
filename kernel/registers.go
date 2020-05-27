// Copyright 2020 Maxim Chernyatevich. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package kernel

const (
	R0 uint16 = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	PC // program counter
	COND
	COUNT
)

// 65536 locations
var Reg [COUNT]uint16

const (
	KBSR uint16 = 0xFE00 // keyboard status
	KBDR uint16 = 0xFE02 // keyboard data
)
