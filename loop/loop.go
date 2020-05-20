package loop

import (
	"fmt"
	"lc3/instructions"
	"lc3/registers"
	"lc3/utils"
)

func Loop() {

	// set the PC to starting position
	// 0x3000 is the default
	var PC_START uint16 = 0x3000

	fmt.Println("Computer starting...")

	registers.Reg[registers.R_PC] = PC_START

	fmt.Println()

	for {
		instruction := utils.MemoryRead(registers.Reg[registers.R_PC])
		opcode := instruction >> 12
		registers.Reg[registers.R_PC]++

		instructions.CallOpcode(opcode, instruction)
	}
}
