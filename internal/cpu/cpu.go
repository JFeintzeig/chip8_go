package cpu

import (
  "bytes"
  "fmt"
  "log"
  "os"
)

var PC uint16 = 0
var I uint16
var DelayTimer uint8
var SoundTimer uint8

var Display = [32*64]uint8{}
// TODO: implement a stack
var Stack [16]uint16
var VariableRegister [16]uint8

var Memory = bytes.NewBuffer(make([]byte, 0, 4096))

var Font = []uint16{
        0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
        0x20, 0x60, 0x20, 0x20, 0x70, // 1
        0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
        0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
        0x90, 0x90, 0xF0, 0x10, 0x10, // 4
        0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
        0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
        0xF0, 0x10, 0x20, 0x40, 0x40, // 7
        0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
        0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
        0xF0, 0x90, 0xF0, 0x90, 0x90, // A
        0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
        0xF0, 0x80, 0x80, 0x80, 0xF0, // C
        0xE0, 0x90, 0x90, 0x90, 0xE0, // D
        0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
        0xF0, 0x80, 0xF0, 0x80, 0x80,  // F
}

// i really want something that allows me to get:
// - the first nibble
// - the second nibble
// - the third nibble
// - the fourth nibble
// - the second, third, fourth together
// - the third, fourth together
// that way i can get the first, route to a instruction function
// via a hashmap based on the first nibble, pass the instruction uint16 (or custom type)
// into the function, and allow the function to parse the other nibbles as needed.

// untested
func getNibble(instruction uint16, i uint8) uint16 {
  if i > 3 {
    log.Fatal("Can't get %d nibble of instruction, it only has 0 - 3.", i)
  }
  return instruction & (0x01 << (3-i))
}

func incrementPC() {
  PC += 2
}

func Fetch(file os.File) uint16 {
  twoBytes := make([]byte, 2)
  offset, err := file.Seek(int64(PC), 0)
  if err != nil {
    log.Fatal("problem seeking in file: %s", err)
  }
  _, err = file.Read(twoBytes)
  if err != nil {
    log.Fatal("problem reading instruction: %s", err)
  }
  instruction := (uint16(twoBytes[0]) << 8) | uint16(twoBytes[1])
  fmt.Printf("two bytes at %d: %x\n", offset, instruction)
  incrementPC()
  return instruction
}

func Decode() {
  fmt.Println("decode")
}

// maybe make hash map of instruction int --> instruction function implementation?
// then Execute just passes 
func Execute(instruction uint16) {
  fmt.Println("execute")
  I00E0()
}

// debugging stuff

func SetDisplay() {
  Display[5] = 1
}

