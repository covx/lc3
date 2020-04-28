package utils

import (
	"bytes"
	"encoding/binary"
	"lc3/conditions"
	"lc3/memory"
	"lc3/registers"
	"log"
	"os"
	"unsafe"
)

var NativeEndian binary.ByteOrder

func SignExtend(number uint16, bitCount int) uint16 {
	if (number >> (bitCount - 1) & 1) != 0 {
		number |= 0xFFFF << bitCount
	}
	return number
}

func UpdateFlags(register uint16) {
	if registers.Reg[register] == 0 {
		registers.Reg[registers.R_COND] = conditions.FL_ZRO
	} else if (registers.Reg[register] >> 15) == 1 {
		// a 1 in the left-most bit indicates negative
		registers.Reg[registers.R_COND] = conditions.FL_NEG
	} else {
		registers.Reg[registers.R_COND] = conditions.FL_POS
	}
}

func swapToLittleEndian16(data uint16) uint16 {
	return (data << 8) | (data >> 8)
}

func readObjFile(path string) ([]byte, int64) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal("file.Stat failed", err)
	}

	var size int64 = info.Size()

	data := make([]byte, size)

	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	return data, size
}

func ReadImageFileToMemory(path string) {
	var header uint16

	data, _ := readObjFile(path)

	buffer := bytes.NewBuffer(data)
	header = binary.BigEndian.Uint16(buffer.Next(2))
	log.Printf("Header has been read: 0x%x", header)

	bufferLen := buffer.Len()
	for i := 0; i < bufferLen; i++ {
		b := buffer.Next(2)
		if len(b) == 0 {
			break
		}
		memory.Memory[i] = binary.BigEndian.Uint16(b)
	}
	log.Printf("Program has been read into memory, contains %d bytes, %d words", bufferLen, bufferLen/2)
}

// nativeEndian is the byte order for the local platform. Used to send back and
// forth Tensors with the C API. We test for endianness at runtime because
// some architectures can be booted into different endian modes.
//	https://github.com/tensorflow/tensorflow/blob/master/tensorflow/go/tensor.go
func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		NativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		NativeEndian = binary.BigEndian
	default:
		log.Fatal("could not determine native endianness")
	}
}
