package kernel

import (
	"fmt"
	"log"
)

func Loop() {

	// set the PC to starting position
	// 0x3000 is the default
	var PCStart uint16 = 0x3000

	log.Println("Computer starting...")

	Register[PC] = PCStart

	fmt.Println()

	for {
		instruction := memory.Read(Register[PC])
		opcode := instruction >> 12
		Register[PC]++

		callOpcode(opcode, instruction)
	}
}
