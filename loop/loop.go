package loop

import (
	"fmt"
	"lc3/instructions"
	"lc3/opcodes"
	"lc3/registers"
	"lc3/utils"
)

func Loop() {

	// set the PC to starting position
	// 0x3000 is the default
	var PC_START uint16 = 0x3000

	fmt.Println("Computer starting...")

	registers.Reg[registers.R_PC] = PC_START

	//var operation uint16

	var currentRegister uint16 = registers.Reg[registers.R_PC]

	for {

		instruction := utils.MemoryRead(currentRegister)
		operation := instruction >> 12
		currentRegister++

		switch operation {
		case opcodes.OP_ADD:
			instructions.Add(instruction)
			fmt.Println("============")
			break
			//case opcodes.OP_AND:
			//	{AND, 7}
			//	break
			//case opcodes.OP_NOT:
			//	{NOT, 7}
			//	break
			//case opcodes.OP_BR:
			//	{BR, 7}
			//	break
			//case opcodes.OP_JMP:
			//	{JMP, 7}
			//	break
			//case opcodes.OP_JSR
			//	{JSR, 7}
			//	break
			//case opcodes.OP_LD
			//	{LD, 7}
			//	break
			//case opcodes.OP_LDI:
			//	{LDI, 6}
			//	break
			//case opcodes.OP_LDR:
			//	{LDR, 7}
			//	break
			//case opcodes.OP_LEA:
			//	{LEA, 7}
			//	break
			//case opcodes.OP_ST:
			//	{ST, 7}
			//	break
			//case opcodes.OP_STI:
			//	{STI, 7}
			//	break
			//case opcodes.OP_STR:
			//	{STR, 7}
			//	break
			//case opcodes.OP_TRAP:
			//	{TRAP, 8}
			//	break
			//case opcodes.OP_RES:
			//	fallthrough
			//case opcodes.OP_RTI:
			//	fallthrough
			//default:
			//	{BAD OPCODE, 7}
			//	break
		}

	}

	// {Shutdown, 12}
	fmt.Println("Computer halting...")
}
