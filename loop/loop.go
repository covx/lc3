package loop

import (
	"fmt"
	"lc3/instructions"
	"lc3/registers"
	"lc3/utils"
	"log"
)

func Loop() {

	// set the PC to starting position
	// 0x3000 is the default
	var PCStart uint16 = 0x3000

	log.Println("Computer starting...")

	registers.Reg[registers.PC] = PCStart

	fmt.Println()

	for {
		instruction := utils.MemoryRead(registers.Reg[registers.PC])
		opcode := instruction >> 12
		registers.Reg[registers.PC]++

		instructions.CallOpcode(opcode, instruction)
	}
}
