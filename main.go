// Copyright 2020 Maxim Chernyatevich. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

// lc3 implements LC3 little computer assembly emulator

package main

import (
	"flag"
	"lc3/kernel"
)

func main() {
	imgFilePath := flag.String("image", "empty image", "go -image program.obj")
	flag.Parse()
	kernel.ReadImageFileToMemory(*imgFilePath)
	kernel.Loop()
}
