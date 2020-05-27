// Copyright 2020 Maxim Chernyatevich. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package kernel

import (
	"log"
)

const (
	BR   uint16 = iota // branch
	ADD                // add
	LD                 // load
	ST                 // store
	JSR                // jump register
	AND                // bitwise and
	LDR                // load register
	STR                // store register
	RTI                // unused
	NOT                // bitwise not
	LDI                // load indirect
	STI                // store indirect
	JMP                // jump
	RES                // reserved (unused)
	LEA                // load effective address
	TRAP               // execute trap
)

var opcodesMapping = map[uint16]func(uint16){
	BR:   branch,
	ADD:  add,
	LD:   load,
	ST:   store,
	JSR:  jumpRegister,
	AND:  and,
	LDR:  loadBaseOffset,
	STR:  storeBaseOffset,
	RTI:  unusedOpcode,
	NOT:  not,
	LDI:  loadIndirect,
	STI:  storeIndirect,
	JMP:  jump,
	RES:  unusedOpcode,
	LEA:  loadEffectiveAddress,
	TRAP: systemCall,
}

func callOpcode(opcode uint16, instruction uint16) {
	if f, ok := opcodesMapping[opcode]; ok {
		f(instruction)
	} else {
		log.Printf("invalid opcode=%v", opcode)
	}

}
