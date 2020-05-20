// Copyright 2020 Maxim Chernyatevich. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

// Package instructions implements bit-wise instructions for lc3 emulator

package instructions

import (
	"lc3/registers"
	"lc3/utils"
	"log"
)

// getR0 gets destination register (DR)
func getR0(instruction uint16) uint16 {
	return (instruction >> 9) & 0x7
}

// getR1 gets first operand (SR1)
func getR1(instruction uint16) uint16 {
	return (instruction >> 6) & 0x7
}

// getImmFlag gets condition flag
func getConditionFlag(instruction uint16) uint16 {
	return (instruction >> 9) & 0x7
}

// Whether we are in immediate mode
func getImmFlag(instruction uint16) uint16 {
	return (instruction >> 5) & 0x1
}

// getPcOffset9 returns a 9-bit offset;
//
// bits [8:0] of an instruction; used with the PC+offset addressing mode.
// Bits [8:0] are taken as a 9-bit signed 2’s complement integer, sign-extended to 16
// bits and then added to the incremented PC to form an address. Range −256..255.
func getPcOffset9(instruction uint16) uint16 {
	return utils.SignExtend(instruction&0x1ff, 9)
}

// add implements addition, ADD;
//
// Assembler Formats:
//
// ADD DR, SR1, SR2
// ADD DR, SR1, imm5
//
// If bit [5] is 0, the second source operand is obtained from SR2. If bit [5] is 1, the
// second source operand is obtained by sign-extending the imm5 field to 16 bits.
// In both cases, the second source operand is added to the contents of SR1 and the
// result stored in DR. The condition codes are set, based on whether the result is
// negative, zero, or positive.
//
// Examples:
//
// ADD R2, R3, R4 ; R2 ← R3 + R4
// ADD R2, R3, #7 ; R2 ← R3 + 7
func add(instruction uint16) {
	var r0 = getR0(instruction)
	var r1 = getR1(instruction)

	if getImmFlag(instruction) == 1 {
		var imm5 = utils.SignExtend(instruction&0x1f, 5)
		registers.Reg[r0] = registers.Reg[r1] + imm5
	} else {
		var r2 = instruction & 0x7
		registers.Reg[r0] = registers.Reg[r1] + registers.Reg[r2]
	}
	utils.UpdateFlags(r0)
}

// and implements bit-wise Logical AND;
//
// Assembler Formats:
//
// AND DR, SR1, SR2
// AND DR, SR1, imm5
//
// If bit [5] is 0, the second source operand is obtained from SR2. If bit [5] is 1,
// the second source operand is obtained by sign-extending the imm5 field to 16
// bits. In either case, the second source operand and the contents of SR1 are bitwise ANDed,
// and the result stored in DR. The condition codes are set, based on
// whether the binary value produced, taken as a 2’s complement integer, is negative,
// zero, or positive.
//
// Examples:
//
// AND R2, R3, R4 ; R2 ← R3 AND R4
// AND R2, R3, #7 ; R2 ← R3 AND 7
func and(instruction uint16) {
	var r0 = getR0(instruction)
	var r1 = getR1(instruction)

	if getImmFlag(instruction) == 1 {
		var imm5 = utils.SignExtend(instruction&0x1f, 5)
		registers.Reg[r0] = registers.Reg[r1] & imm5
	} else {
		var r2 = instruction & 0x7
		registers.Reg[r0] = registers.Reg[r1] & registers.Reg[r2]
	}
	utils.UpdateFlags(r0)

}

func not(instruction uint16) {
	r0 := getR0(instruction)

	registers.Reg[getR0(instruction)] = ^registers.Reg[getR1(instruction)]
	utils.UpdateFlags(r0)
}

func branch(instruction uint16) {
	if getConditionFlag(instruction)&registers.Reg[registers.R_COND] != 0 {
		registers.Reg[registers.R_PC] += getPcOffset9(instruction)
	}

}

func jump(instruction uint16) {
	registers.Reg[registers.R_PC] = registers.Reg[getR1(instruction)]
}

func jumpRegister(instruction uint16) {
	longPcOffset := utils.SignExtend(instruction&0x7ff, 11)
	longFlag := (instruction >> 11) & 0x1

	registers.Reg[registers.R_R7] = registers.Reg[registers.R_PC]

	if longFlag == 1 {
		registers.Reg[registers.R_PC] += longPcOffset // JSR
	} else {
		registers.Reg[registers.R_PC] = registers.Reg[getR1(instruction)] // JSRR
	}
}

func load(instruction uint16) {
	r0 := getR0(instruction)

	registers.Reg[r0] = utils.MemoryRead(registers.Reg[registers.R_PC] + getPcOffset9(instruction))
	utils.UpdateFlags(r0)

}

func loadIndirect(instruction uint16) {
	r0 := getR0(instruction)

	// Add pcOffset to the current PC, look at that memory location to get the final address
	registers.Reg[r0] = utils.MemoryRead(
		utils.MemoryRead(registers.Reg[registers.R_PC] + getPcOffset9(instruction)))
	utils.UpdateFlags(r0)
}

func loadBaseOffset(instruction uint16) {
	r0 := getR0(instruction)
	r1 := getR1(instruction)
	offset := utils.SignExtend(instruction&0x3f, 6)

	registers.Reg[r0] = utils.MemoryRead(registers.Reg[r1] + offset)
	utils.UpdateFlags(r0)
}

func loadEffectiveAddress(instruction uint16) {
	r0 := getR0(instruction)

	registers.Reg[r0] = registers.Reg[registers.R_PC] + getPcOffset9(instruction)
	utils.UpdateFlags(r0)
}

func store(instruction uint16) {
	utils.MemoryWrite(
		registers.Reg[registers.R_PC]+getPcOffset9(instruction),
		registers.Reg[getR0(instruction)])
}

func storeIndirect(instruction uint16) {
	utils.MemoryWrite(
		utils.MemoryRead(registers.Reg[registers.R_PC]+getPcOffset9(instruction)),
		registers.Reg[getR0(instruction)])
}

func storeBaseOffset(instruction uint16) {
	offset := utils.SignExtend(instruction&0x3f, 6)

	utils.MemoryWrite(
		registers.Reg[getR1(instruction)]+offset, registers.Reg[getR0(instruction)])
}

func unusedOpcode(instruction uint16) {
	log.Printf("invalid insruction=%v", instruction)
}
