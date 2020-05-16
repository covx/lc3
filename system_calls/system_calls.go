package system_calls

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"lc3/memory"
	"lc3/opcodes"
	"lc3/registers"
	"log"
	"os"
)

func KeyboardRead() uint16 {
	var keyPressed uint16 = 0x0

	symb, controlKey, err := keyboard.GetSingleKey()
	keyPressed = uint16(symb)

	if controlKey == keyboard.KeyEsc || controlKey == keyboard.KeyCtrlC {
		fmt.Println("\n\nPressed escaping")
		haltComputer()
	}

	if err != nil {
		log.Printf("Error, %s", err)
	}
	return keyPressed
}

// Reads a single ASCII char; Trap GETC
func readCharFromKeyboard() {
	registers.Reg[registers.R_R0] = KeyboardRead()
}

// Prints a single character to sdtout
func outCharToStdout() {
	fmt.Printf("%c\n", registers.Reg[registers.R_R0])
}

// Halts computer; breaks main loop
func haltComputer() {
	fmt.Println("Computer halting...")
	os.Exit(0)
}

// Writes a string of ASCII characters to the console display
func outStringToStdout() {
	for address := registers.Reg[registers.R_R0]; memory.Memory[address] != 0x00; address++ {
		fmt.Printf("%c", memory.Memory[address])
	}
}

func SystemCall(instruction uint16) {

	switch instruction & 0xff {
	case opcodes.TRAP_GETC:
		readCharFromKeyboard()
		break
	case opcodes.TRAP_OUT:
		outCharToStdout()
		break
	case opcodes.TRAP_PUTS:
		outStringToStdout()
		break
	case opcodes.TRAP_IN:
		//{TRAP IN, 9}
		fmt.Println("TRAPIN")
		break
	case opcodes.TRAP_PUTSP:
		fmt.Println("TRAPPUTSP")
		//{TRAP PUTSP, 9}
		break
	case opcodes.TRAP_HALT:
		haltComputer()
		break
	}
}
