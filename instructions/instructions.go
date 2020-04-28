package instructions

import (
	"lc3/registers"
	"lc3/utils"
)

func Add(instruction uint16) {
	// destination register (DR)
	var r0 = (instruction >> 9) & 0x7

	// first operand (SR1)
	var r1 = (instruction >> 6) & 0x7

	// whether we are in immediate mode
	var immFlag = (instruction >> 5) & 0x1

	if immFlag == 1 {
		var imm5 = utils.SignExtend(instruction&0x1F, 5)
		registers.Reg[r0] = registers.Reg[r1] + imm5
	} else {
		var r2 = instruction & 0x7
		registers.Reg[r0] = registers.Reg[r1] + registers.Reg[r2]
	}
	utils.UpdateFlags(r0)
}
