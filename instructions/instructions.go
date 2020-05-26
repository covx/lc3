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

// not implements bit-wise Bit-Wise Complement, NOT;
//
// Assembler Formats:
//
// NOT DR, SR
//
// The bit-wise complement of the contents of SR is stored in DR.
// The condition codes are set, based on whether the binary value produced, taken as a 2’s
// complement integer, is negative, zero, or positive.
//
// Examples:
//
// NOT R4, R2 ; R4 ← NOT(R2)
func not(instruction uint16) {
	r0 := getR0(instruction)

	registers.Reg[getR0(instruction)] = ^registers.Reg[getR1(instruction)]
	utils.UpdateFlags(r0)
}

// branch implements Conditional Branch, BR;
//
// Assembler Formats:
//
// BRn LABEL BRzp LABEL
// BRz LABEL BRnp LABEL
// BRp LABEL BRnz LABEL
//
// The condition codes specified by the state of bits [11:9] are tested. If bit [11] is
// set, N is tested; if bit [11] is clear, N is not tested. If bit [10] is set, Z is tested, etc.
// If any of the condition codes tested is set, the program branches to the location
// specified by adding the sign-extended PCoffset9 field to the incremented PC.
//
// Examples:
//
// BRzp LOOP ; Branch to LOOP if the last result was zero or positive.
func branch(instruction uint16) {
	if getConditionFlag(instruction)&registers.Reg[registers.R_COND] != 0 {
		registers.Reg[registers.R_PC] += getPcOffset9(instruction)
	}

}

// jump implements Jump, JMP;
//
// Assembler Formats:
//
// JMP BaseR
//
// The program unconditionally jumps to the location specified by the contents of
// the base register. Bits [8:6] identify the base register.
//
// Examples:
//
// JMP R2 ; PC ← R2
func jump(instruction uint16) {
	registers.Reg[registers.R_PC] = registers.Reg[getR1(instruction)]
}

// jumpRegister implements Jump to Subroutine -- JSR, JSRR;
//
// Assembler Formats:
//
// JSR LABEL
// JSRR BaseR
//
// First, the incremented PC is saved in R7. This is the linkage back to the calling
// routine. Then the PC is loaded with the address of the first instruction of the
// subroutine, causing an unconditional jump to that address. The address of the
// subroutine is obtained from the base register (if bit [11] is 0), or the address is
// computed by sign-extending bits [10:0] and adding this value to the incremented
// PC (if bit [11] is 1).
//
// Examples:
//
// JSR QUEUE ; Put the address of the instruction following JSR into R7;
//			 ; Jump to QUEUE.
// JSRR R3   ; Put the address following JSRR into R7; Jump to the
//			 ; address contained in R3.
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

// load implements Load, LD;
//
// Assembler Formats:
//
// LD DR, LABEL
//
// An address is computed by sign-extending bits [8:0] to 16 bits and adding this
// value to the incremented PC. The contents of memory at this address are loaded
// into DR. The condition codes are set, based on whether the value loaded is
// negative, zero, or positive.
//
// Examples:
//
// LD R4, VALUE ; R4 ← mem[VALUE]
func load(instruction uint16) {
	r0 := getR0(instruction)

	registers.Reg[r0] = utils.MemoryRead(registers.Reg[registers.R_PC] + getPcOffset9(instruction))
	utils.UpdateFlags(r0)

}

// loadIndirect implements Load Indirect, LDI;
//
// Assembler Formats:
//
// LDI DR, LABEL
//
// An address is computed by sign-extending bits [8:0] to 16 bits and adding this
// value to the incremented PC. What is stored in memory at this address is the
// address of the data to be loaded into DR. The condition codes are set, based on
// whether the value loaded is negative, zero, or positive.
//
// Examples:
//
// LDI R4, ONEMORE ; R4 ← mem[mem[ONEMORE]]
func loadIndirect(instruction uint16) {
	r0 := getR0(instruction)

	// Add pcOffset to the current PC, look at that memory location to get the final address
	registers.Reg[r0] = utils.MemoryRead(
		utils.MemoryRead(registers.Reg[registers.R_PC] + getPcOffset9(instruction)))
	utils.UpdateFlags(r0)
}

// loadBaseOffset implements Load Base+offset, LDR;
//
// Assembler Formats:
//
// LDR DR, BaseR, offset6
//
// An address is computed by sign-extending bits [5:0] to 16 bits and adding this
// value to the contents of the register specified by bits [8:6]. The contents of memory
// at this address are loaded into DR. The condition codes are set, based on whether
// the value loaded is negative, zero, or positive
//
// Examples:
//
// LDR R4, R2, #−5 ; R4 ← mem[R2 − 5]
func loadBaseOffset(instruction uint16) {
	r0 := getR0(instruction)
	r1 := getR1(instruction)
	offset := utils.SignExtend(instruction&0x3f, 6)

	registers.Reg[r0] = utils.MemoryRead(registers.Reg[r1] + offset)
	utils.UpdateFlags(r0)
}

// loadEffectiveAddress implements Load Effective Address, LEA;
//
// Assembler Formats:
//
// LEA DR, LABEL
//
// An address is computed by sign-extending bits [8:0] to 16 bits and adding this
// value to the incremented PC. This address is loaded into DR.‡ The condition
// codes are set, based on whether the value loaded is negative, zero, or positive.
//
// Examples:
//
// LEA R4, TARGET ; R4 ← address of TARGET.
func loadEffectiveAddress(instruction uint16) {
	r0 := getR0(instruction)

	registers.Reg[r0] = registers.Reg[registers.R_PC] + getPcOffset9(instruction)
	utils.UpdateFlags(r0)
}

// store implements Store, ST;
//
// Assembler Formats:
//
// ST SR, LABEL
//
// The contents of the register specified by SR are stored in the memory location
// whose address is computed by sign-extending bits [8:0] to 16 bits and adding this
// value to the incremented PC.
//
// Examples:
//
// ST R4, HERE ; mem[HERE] ← R4
func store(instruction uint16) {
	utils.MemoryWrite(
		registers.Reg[registers.R_PC]+getPcOffset9(instruction),
		registers.Reg[getR0(instruction)])
}

// storeIndirect implements Store Indirect, STI;
//
// Assembler Formats:
//
// STI SR, LABEL
//
// The contents of the register specified by SR are stored in the memory location
// whose address is obtained as follows: Bits [8:0] are sign-extended to 16 bits and
// added to the incremented PC. What is in memory at this address is the address of
// the location to which the data in SR is stored.
//
// Examples:
//
// STI R4, NOT_HERE ; mem[mem[NOT_HERE]] ← R4
func storeIndirect(instruction uint16) {
	utils.MemoryWrite(
		utils.MemoryRead(registers.Reg[registers.R_PC]+getPcOffset9(instruction)),
		registers.Reg[getR0(instruction)])
}

// storeBaseOffset implements Store Base+offset, STR;
//
// Assembler Formats:
//
// STR SR, BaseR, offset6
//
// The contents of the register specified by SR are stored in the memory location
// whose address is computed by sign-extending bits [5:0] to 16 bits and adding this
// value to the contents of the register specified by bits [8:6].
//
// Examples:
//
// STR R4, R2, #5 ; mem[R2 + 5] ← R4
func storeBaseOffset(instruction uint16) {
	offset := utils.SignExtend(instruction&0x3f, 6)

	utils.MemoryWrite(
		registers.Reg[getR1(instruction)]+offset, registers.Reg[getR0(instruction)])
}

// unusedOpcode prints warning into stdout if instruction isn`t valid;
func unusedOpcode(instruction uint16) {
	log.Printf("invalid insruction=%v", instruction)
}
