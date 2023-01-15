package cpu

import (
  "fmt"
  "log"
  "os"
)

type instruction struct {
  full uint16
  a uint8
  x uint8
  y uint8
  n uint8
  nn uint8
  nnn uint16
}

type chip8 struct {
  pc uint16
  i uint16
  delayTimer uint8
  soundTimer uint8
  display [32*64]uint8
  // TODO: implement a stack
  stack [16]uint16
  variableRegister [16]uint8
  memory [4096]byte
  instructionMap map[uint8]func(*instruction)
}

func (c8 *chip8) incrementPC() {
  c8.pc += 2
}

// TODO: split this up and add functionality
//   1. load entire file into `memory` starting at address 0x200
//   2. fetch and decode function that accesses that memory directly
//   need to set starting bound on where to place in memory, abstract it away
func (c8 *chip8) FetchAndDecode(file os.File) instruction {
  twoBytes := make([]byte, 2)
  offset, err := file.Seek(int64(c8.pc), 0)
  if err != nil {
    log.Fatal("problem seeking in file: %s", err)
  }
  _, err = file.Read(twoBytes)
  if err != nil {
    log.Fatal("problem reading instruction: %s", err)
  }
  coded_instruction := (uint16(twoBytes[0]) << 8) | uint16(twoBytes[1])
  fmt.Printf("two bytes at %d: %x\n", offset, coded_instruction)
  c8.incrementPC()
  return instruction{
    coded_instruction,
    uint8((coded_instruction & 0xF000) >> 12),
    uint8((coded_instruction & 0x0F00) >> 8),
    uint8((coded_instruction & 0x00F0) >> 4),
    uint8(coded_instruction & 0x000F),
    uint8(coded_instruction & 0x00FF),
    coded_instruction & 0x0FFF,
  }
}

func (c8 *chip8) Execute(instruction *instruction) {
  fmt.Printf("A, X, Y, N, NN, NNN: %x, %x, %x, %x, %x, %x\n", instruction.a, instruction.x, instruction.y, instruction.n, instruction.nn, instruction.nnn)
  if instructionFunc, ok := c8.instructionMap[instruction.a]; ok {
    instructionFunc(instruction)
  } else {
    log.Fatalf("no instruction for %x, first nibble %x", instruction.full, instruction.a)
  }
}

var font = []uint8{
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

func NewChip8() *chip8 {
  // load font into memory starting at 0x050
  // TODO: abstract away this memory starting place
  memory := [4096]byte{}

  for index, element := range font {
    memory[0x050 + index] = element
  }

  instructionMap := map[uint8]func(*instruction){}

  c8 := chip8{0, 0, 0, 0, [32*64]uint8{}, [16]uint16{}, [16]uint8{}, memory, instructionMap}

  // put instructions in a map
  c8.instructionMap[0x0] = c8.I00E0
  c8.instructionMap[0x1] = c8.I1NNN
  c8.instructionMap[0x6] = c8.I6XNN
  c8.instructionMap[0x7] = c8.I7XNN
  c8.instructionMap[0xA] = c8.IANNN

  return &c8
}

// debugging stuff

func (c8 *chip8) SetDisplay() {
  c8.display[5] = 1
}
