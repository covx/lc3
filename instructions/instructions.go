package instructions

import (
	"lc3/registers"
	"lc3/utils"
	"log"
)

func add(instruction uint16) {
	// destination register (DR)
	var r0 = (instruction >> 9) & 0x7

	// first operand (SR1)
	var r1 = (instruction >> 6) & 0x7

	// whether we are in immediate mode
	var immFlag = (instruction >> 5) & 0x1

	if immFlag == 1 {
		var imm5 = utils.SignExtend(instruction&0x1f, 5)
		registers.Reg[r0] = registers.Reg[r1] + imm5
	} else {
		var r2 = instruction & 0x7
		registers.Reg[r0] = registers.Reg[r1] + registers.Reg[r2]
	}
	utils.UpdateFlags(r0)
}

func and(instruction uint16) {
	var r0 = (instruction >> 9) & 0x7
	var r1 = (instruction >> 6) & 0x7
	var immFlag = (instruction >> 5) & 0x1

	if immFlag == 1 {
		var imm5 = utils.SignExtend(instruction&0x1f, 5)
		registers.Reg[r0] = registers.Reg[r1] & imm5
	} else {
		var r2 = instruction & 0x7
		registers.Reg[r0] = registers.Reg[r1] & registers.Reg[r2]
	}
	utils.UpdateFlags(r0)

}

func not(instruction uint16) {
	var r0 = (instruction >> 9) & 0x7
	var r1 = (instruction >> 6) & 0x7

	registers.Reg[r0] = ^registers.Reg[r1]
	utils.UpdateFlags(r0)
}

func branch(instruction uint16) {
	var pcOffset = utils.SignExtend(instruction&0x1ff, 9)
	var conditionFlag = (instruction >> 9) & 0x7

	if conditionFlag&registers.Reg[registers.R_COND] != 0 {
		registers.Reg[registers.R_PC] += pcOffset
	}

}

func jump(instruction uint16) {
	// Also handles RET
	var r1 = (instruction >> 6) & 0x7

	registers.Reg[registers.R_PC] = registers.Reg[r1]
}

func jumpRegister(instruction uint16) {
	var r1 = (instruction >> 6) & 0x7
	var longPcOffset = utils.SignExtend(instruction&0x7ff, 11)
	var longFlag = (instruction >> 11) & 0x1

	registers.Reg[registers.R_R7] = registers.Reg[registers.R_PC]

	if longFlag == 1 {
		registers.Reg[registers.R_PC] += longPcOffset // JSR
	} else {
		registers.Reg[registers.R_PC] = registers.Reg[r1] // JSRR
	}
}

func load(instruction uint16) {
	var r0 = (instruction >> 9) & 0x7
	var pcOffset = utils.SignExtend(instruction&0x1ff, 9)

	registers.Reg[r0] = utils.MemoryRead(registers.Reg[registers.R_PC] + pcOffset)
	utils.UpdateFlags(r0)

}

func loadIndirect(instruction uint16) {
	var r0 = (instruction >> 9) & 0x7
	var pcOffset = utils.SignExtend(instruction&0x1ff, 9)

	// Add pcOffset to the current PC, look at that memory location to get the final address
	registers.Reg[r0] = utils.MemoryRead(utils.MemoryRead(registers.Reg[registers.R_PC] + pcOffset))
	utils.UpdateFlags(r0)
}

func loadBaseOffset(instruction uint16) {
	var r0 = (instruction >> 9) & 0x7
	var r1 = (instruction >> 6) & 0x7

	var offset = utils.SignExtend(instruction&0x3f, 6)

	registers.Reg[r0] = utils.MemoryRead(registers.Reg[r1] + offset)
	utils.UpdateFlags(r0)
}

func loadEffectiveAddress(instruction uint16) {
	var r0 = (instruction >> 9) & 0x7
	var pcOffset = utils.SignExtend(instruction&0x1ff, 9)

	registers.Reg[r0] = registers.Reg[registers.R_PC] + pcOffset
	utils.UpdateFlags(r0)
}

func store(instruction uint16) {
	var r0 = (instruction >> 9) & 0x7
	var pcOffset = utils.SignExtend(instruction&0x1ff, 9)

	utils.MemoryWrite(registers.Reg[registers.R_PC]+pcOffset, registers.Reg[r0])
}

func storeIndirect(instruction uint16) {
	var r0 = (instruction >> 9) & 0x7
	var pcOffset = utils.SignExtend(instruction&0x1ff, 9)

	utils.MemoryWrite(utils.MemoryRead(registers.Reg[registers.R_PC]+pcOffset), registers.Reg[r0])
}

func storeBaseOffset(instruction uint16) {
	var r0 = (instruction >> 9) & 0x7
	var r1 = (instruction >> 6) & 0x7
	var offset = utils.SignExtend(instruction&0x3f, 6)

	utils.MemoryWrite(registers.Reg[r1]+offset, registers.Reg[r0])
}

func unusedOpcode(instruction uint16) {
	log.Printf("invalid insruction=%v", instruction)
}
