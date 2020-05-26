// Copyright 2020 Maxim Chernyatevich. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

// Package system_calls implements TRAP instructions for lc3 emulator

package system_calls

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"lc3/memory"
	"lc3/registers"
	"log"
	"os"
)

const (
	GETC  uint16 = 0x20 // get character from keyboard
	OUT   uint16 = 0x21 // output a character
	PUTS  uint16 = 0x22 // output a word string
	IN    uint16 = 0x23 // input a string
	PUTSP uint16 = 0x24 // output a byte string
	HALT  uint16 = 0x25 // halt the program
)

// Reads single symbol from keyboard
func KeyboardRead() uint16 {
	symb, controlKey, err := keyboard.GetSingleKey()

	if controlKey == keyboard.KeyEsc || controlKey == keyboard.KeyCtrlC {
		log.Println("Pressed escaping")
		haltComputer()
	}

	if err != nil {
		log.Printf("Error, %s", err)
	}
	return uint16(symb)
}

// Reads a single ASCII char; Trap GETC
func readCharFromKeyboard() {
	registers.Reg[registers.R0] = KeyboardRead()
}

// Prints a single character to stdout
func outCharToStdout() {
	fmt.Printf("%c", registers.Reg[registers.R0])
}

// Halts computer; breaks main loop
func haltComputer() {
	log.Printf("Computer halting...")
	os.Exit(0)
}

// Writes a string of ASCII characters to the console display
func outStringToStdout() {
	for address := registers.Reg[registers.R0]; memory.Memory[address] != 0x00; address++ {
		fmt.Printf("%c", memory.Memory[address])
	}
}

// Prints a prompt on the screen and
//reads a single character from the keyboard
func printPromtAndRead() {
	fmt.Printf("Enter a character: ")
	symb := KeyboardRead()
	fmt.Printf("%c", symb)
	registers.Reg[registers.R0] = symb
}

// Write a string of ASCII characters to the stdout
func printStringToConsole() {
	for address := registers.Reg[registers.R0]; memory.Memory[address] != 0x00; address++ {
		value := memory.Memory[address]

		fmt.Printf("%c", value&0xff)

		symb := value & 0xff >> 8
		if symb != 0 {
			fmt.Printf("%c", symb)
		}
	}
}

var callMapping = map[uint16]func(){
	GETC:  readCharFromKeyboard,
	OUT:   outCharToStdout,
	PUTS:  outStringToStdout,
	IN:    printPromtAndRead,
	PUTSP: printStringToConsole,
	HALT:  haltComputer,
}

func SystemCall(instruction uint16) {
	callMapping[instruction&0xff]()
}
