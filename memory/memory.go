// Copyright 2020 Maxim Chernyatevich. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

// Package memory implements memory for lc3 emulator

package memory

import "math"

// 65536 locations
var Memory [math.MaxUint16 + 1]uint16
