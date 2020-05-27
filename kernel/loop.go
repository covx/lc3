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

	Reg[PC] = PCStart

	fmt.Println()

	for {
		instruction := memory.Read(Reg[PC])
		opcode := instruction >> 12
		Reg[PC]++

		callOpcode(opcode, instruction)
	}
}
