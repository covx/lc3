// Copyright 2020 Maxim Chernyatevich. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package kernel

import (
	"math"
)

// 65536 locations
type lc3Memory [math.MaxUint16 + 1]uint16

func (m *lc3Memory) Write(address uint16, val uint16) {
	m[address] = val
}

func (m *lc3Memory) Read(address uint16) uint16 {
	if address == KBSR {
		checkKey := keyboardRead()

		if checkKey != 0 {
			m[KBSR] = 1 << 15
			m[KBDR] = checkKey
		} else {
			m[KBSR] = 0
		}
	}
	return m[address]
}

var memory lc3Memory
