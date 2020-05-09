package main

import (
	"flag"
	"lc3/loop"
	"lc3/utils"
)

func main() {
	imgFilePath := flag.String("image", "empty image", "go -image program.obj")
	flag.Parse()
	utils.ReadImageFileToMemory(*imgFilePath)
	//utils.KeyboardRead()
	loop.Loop()
}
