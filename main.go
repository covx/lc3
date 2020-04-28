package main

import (
	"flag"
	"lc3/utils"
)

func main() {
	imgFilePath := flag.String("image", "empty image", "go -image")
	flag.Parse()
	utils.ReadImageFileToMemory(*imgFilePath)

}
