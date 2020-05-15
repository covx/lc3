package loop

import (
	"fmt"
	"lc3/instructions"
	"lc3/opcodes"
	"lc3/registers"
	"lc3/system_calls"
	"lc3/utils"
	"log"
)

func Loop() {

	// set the PC to starting position
	// 0x3000 is the default
	var PC_START uint16 = 0x3000

	fmt.Println("Computer starting...")

	registers.Reg[registers.R_PC] = PC_START

	for running := 1; running != 0; {
		instruction := utils.MemoryRead(registers.Reg[registers.R_PC])
		operation := instruction >> 12

		registers.Reg[registers.R_PC]++

		switch operation {
		case opcodes.OP_ADD:
			instructions.Add(instruction)
			break
		case opcodes.OP_AND:
			instructions.And(instruction)
			break
		case opcodes.OP_NOT:
			instructions.Not(instruction)
			break
		case opcodes.OP_BR:
			instructions.Branch(instruction)
			break
		case opcodes.OP_JMP:
			instructions.Jump(instruction)
			break
		case opcodes.OP_JSR:
			instructions.JumpRegister(instruction)
			break
		case opcodes.OP_LD:
			instructions.Load(instruction)
			break
		case opcodes.OP_LDI:
			instructions.LoadIndirect(instruction)
			break
		case opcodes.OP_LDR:
			instructions.LoadBaseOffset(instruction)
			break
		case opcodes.OP_LEA:
			instructions.LoadEffectiveAddress(instruction)
			break
		case opcodes.OP_ST:
			instructions.Store(instruction)
			break
		case opcodes.OP_STI:
			instructions.StoreIndirect(instruction)
			break
		case opcodes.OP_STR:
			instructions.StoreBaseOffset(instruction)
			break
		case opcodes.OP_TRAP:
			running = system_calls.SystemCall(instruction)
			break
		case opcodes.OP_RES:
			fallthrough
		case opcodes.OP_RTI:
			fallthrough
		default:
			log.Printf("invalid opcode=%v", operation)
			break
		}

	}

	// {Shutdown, 12}
	fmt.Println("Computer halting...")
}
