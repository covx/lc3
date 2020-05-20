package instructions

import (
	"lc3/system_calls"
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
	TRAP: system_calls.SystemCall,
}

func CallOpcode(opcode uint16, instruction uint16) {
	if f, ok := opcodesMapping[opcode]; ok {
		f(instruction)
	} else {
		log.Printf("invalid opcode=%v", opcode)
	}

}
