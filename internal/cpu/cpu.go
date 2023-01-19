package cpu

import (
  "bufio"
  "errors"
  "fmt"
  "io"
  "log"
  "os"
  "time"
)

const PROGRAM_START uint16 = 0x200
const MAX_PROGRAM_ADDRESS uint16 = 0xE8F
const FONT_START uint16 = 0x050
const CLOCK_SPEED uint16 = 500

type instruction struct {
  full uint16
  a uint8
  x uint8
  y uint8
  n uint8
  nn uint8
  nnn uint16
}

type DebugState int

const (
  PAUSED DebugState = iota
  RUNNING
)

// TODO: i had to export Chip8 to use it in Display.go...another way?
type Chip8 struct {
  pc uint16
  i uint16
  delayTimer uint8
  soundTimer uint8
  Display [32*64]uint8
  // TODO: implement a stack
  stack [16]uint16
  variableRegister [16]uint8
  memory [4096]byte
  instructionMap map[uint8]func(*instruction)
  clockSpeed uint16
  // TODO: think where debug should live / interact
  // could maybe have a debugger struct w/methods that interfaces with Chip8...if part
  // of package it can see private vars...?
  modern bool
  debug bool
  debugState DebugState
  debugBreakpoint uint16
}

func (c8 *Chip8) incrementPC() {
  c8.pc += 2
}

func (c8 *Chip8) LoadFile(filePath string) {
  file, err := os.Open(filePath)
  if err != nil {
    log.Fatal("can't find file")
  }
  br := bufio.NewReader(file)
  i := uint16(0)
  for {
    b, err := br.ReadByte()
    if err != nil && !errors.Is(err, io.EOF) {
        fmt.Println(err)
        break
    }
    c8.memory[PROGRAM_START + i] = b
    if err != nil {
        // end of file
        break
    }
    if i > MAX_PROGRAM_ADDRESS {
      log.Fatalf("Programs can only write between %x and %x in memory, loading this one overflowed",PROGRAM_START,MAX_PROGRAM_ADDRESS)
    }
    i++
  }
}

// TODO: error handling that we don't outstep memory
func (c8 *Chip8) fetchAndDecode() instruction {
  twoBytes := c8.memory[c8.pc:c8.pc+2]
  coded_instruction := (uint16(twoBytes[0]) << 8) | uint16(twoBytes[1])
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

func (c8 *Chip8) executeInstruction(instruction *instruction) {
  if instructionFunc, ok := c8.instructionMap[instruction.a]; ok {
    instructionFunc(instruction)
  } else {
    log.Fatalf("no instruction for %x, first nibble %x", instruction.full, instruction.a)
  }
}

func (c8 *Chip8) debugInstruction(instruction *instruction) {
  c8.executeInstruction(instruction)
  if c8.pc == c8.debugBreakpoint {
    c8.debugState = PAUSED
  }
  DebugLoop:
    for c8.debugState == PAUSED {
      fmt.Printf("PC (incremented): %x; Instruction just executed: %x; A: %x; X: %x; Y: %x; N: %x; NN: %x; NNN: %x\n",
          c8.pc, instruction.full, instruction.a, instruction.x, instruction.y, instruction.n, instruction.nn, instruction.nnn)
      fmt.Printf("Debug: (s)tate, (c)ontinue, (n)ext, (q)uit, (b)reakpoint, $<variable>, (h)elp\n")
      command, _ := bufio.NewReader(os.Stdin).ReadString('\n')
      switch string(command[0]) {
      case "s":
        // TODO
        c8.prettyPrint()
      case "c":
        c8.debugState = RUNNING
      case "n":
        break DebugLoop
      case "q":
        os.Exit(0)
      case "b":
        // TODO: parsing doesn't work b/c i read a string i think...
        c8.debugBreakpoint = uint16(command[0]) << 8 | uint16(command[1])
        fmt.Printf("%x", c8.debugBreakpoint)
        fmt.Printf("%s", command)
        c8.debugState = RUNNING
      case "$":
        // TODO
        c8.safeValuePrint()
      case "h":
        fmt.Println(
          `Commands:
           > s    Print entire Chip8 state to console
           > c    Continue with execution without stopping
           > n    Execute instruction then pause again
           > q    Quit progrm
           > b    Set a breakpoint as a uint16 memory address in hex
           > $    View a Chip8 field, e.g. $variableRegister[2]
           > h    Help, print this message
          `)
      default:
        fmt.Printf("Sorry, %s is not a validcommand", command)
      }
    }
}

func (c8 *Chip8) prettyPrint() {
}

func (c8 *Chip8) safeValuePrint() {
}

func (c8 *Chip8) Execute() {
  for {
    instruction := c8.fetchAndDecode()
    if c8.debug {
      c8.debugInstruction(&instruction)
    } else {
      c8.executeInstruction(&instruction)
    }
    time.Sleep(time.Duration(1000/c8.clockSpeed) * time.Millisecond)
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

func NewChip8(debug bool, modern bool) *Chip8 {
  // load font into memory starting at FONT_START
  memory := [4096]byte{}

  for index, element := range font {
    memory[FONT_START + uint16(index)] = element
  }

  instructionMap := map[uint8]func(*instruction){}

  var debugState DebugState
  if debug {
    debugState = PAUSED
  } else {
    debugState = RUNNING
  }

  debugBreakpoint := uint16(0x0000)

  c8 := Chip8{PROGRAM_START, 0, 0, 0, [32*64]uint8{}, [16]uint16{}, [16]uint8{}, memory, instructionMap, CLOCK_SPEED, modern, debug, debugState, debugBreakpoint}

  // put instructions in a map
  c8.instructionMap[0x0] = c8.I00E0
  c8.instructionMap[0x1] = c8.I1NNN
  c8.instructionMap[0x6] = c8.I6XNN
  c8.instructionMap[0x7] = c8.I7XNN
  c8.instructionMap[0xA] = c8.IANNN
  c8.instructionMap[0xD] = c8.IDXYN

  return &c8
}
