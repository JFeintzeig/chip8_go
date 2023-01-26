package cpu

import (
  "bufio"
  "encoding/binary"
  "encoding/hex"
  "errors"
  "fmt"
  "io"
  "log"
  "os"
  "time"

  "jfeintzeig/chip8/internal/utils"
)

const PROGRAM_START uint16 = 0x200
const MAX_PROGRAM_ADDRESS uint16 = 0xE8F
const FONT_START uint16 = 0x050
const CLOCK_SPEED uint16 = 500
const DELAY_SOUND_TIMER_UPDATE uint16 = 60

type DebugState int

const (
  PAUSED DebugState = iota
  RUNNING
)

type keypress struct {
  Pressed bool
  JustReleased bool
}

func NewKeypress() keypress {
  return keypress{false,false}
}

// TODO: i had to export Chip8 to use it in Display.go...another way?
type Chip8 struct {
  pc uint16
  i uint16
  delayTimer uint8
  soundTimer uint8
  Display [32*64]uint8
  stack utils.Stack
  variableRegister [16]uint8
  memory [4096]byte
  instructionMap map[uint8]func(*utils.Instruction)
  Keyboard *[16]keypress
  clockSpeed uint16
  modern bool
  // TODO: refactor debugger into separate struct in same package
  // could maybe have a debugger struct w/methods that interfaces with Chip8...if part
  // of package it can see private vars...?
  // TODO: figure out how to get debugger to work when program requires key presses
  debug bool
  debugState DebugState
  debugBreakpoint uint16
}

func (c8 *Chip8) GetSoundTimer() uint8 {
  return c8.soundTimer
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
func (c8 *Chip8) fetchAndDecode() utils.Instruction {
  twoBytes := c8.memory[c8.pc:c8.pc+2]
  codedInstruction := (uint16(twoBytes[0]) << 8) | uint16(twoBytes[1])
  c8.incrementPC()
  return utils.InstructionFromBytecode(codedInstruction)
}

func (c8 *Chip8) executeInstruction(instruction *utils.Instruction) {
  if instructionFunc, ok := c8.instructionMap[instruction.A]; ok {
    instructionFunc(instruction)
  } else {
    log.Fatalf("no instruction for %x, first nibble %x", instruction.Full, instruction.A)
  }
}

func (c8 *Chip8) debugInstruction(instruction *utils.Instruction) {
  c8.executeInstruction(instruction)
  if c8.pc == c8.debugBreakpoint {
    c8.debugState = PAUSED
  }
  DebugLoop:
    for c8.debugState == PAUSED {
      fmt.Printf("PC (incremented): %x; Instruction just executed: %x; A: %x; X: %x; Y: %x; N: %x; NN: %x; NNN: %x\n",
          c8.pc, instruction.Full, instruction.A, instruction.X, instruction.Y, instruction.N, instruction.NN, instruction.NNN)
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
        fmt.Printf("Enter memory address in hex (e.g. 0xXXXX): ")
        memToBreak, _ := bufio.NewReader(os.Stdin).ReadString('\n')
        val, _ := hex.DecodeString(string(memToBreak[2:]))
        c8.debugBreakpoint = binary.BigEndian.Uint16(val)
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
           > q    Quit program
           > b <0xXXXX>   Set a breakpoint as a uint16 memory address in hex
           > $    View a Chip8 field, e.g. $variableRegister[2]
           > h    Help, print this message
          `)
      default:
        fmt.Printf("Sorry, %s is not a validcommand", command)
      }
    }
}

func (c8 *Chip8) prettyPrint() {
  fmt.Printf("%x\n",c8.pc)
  fmt.Printf("%x\n",c8.i)
  fmt.Printf("%x\n", c8.variableRegister)
}

func (c8 *Chip8) safeValuePrint() {
}

// TODO: run some more test games, especially
// newer stuff from https://johnearnest.github.io/chip8Archive/?sort=platform
// see if i can find other bugs in my code
// write unit tests
func (c8 *Chip8) Execute() {
  loopCounter := 0
  for {
    instruction := c8.fetchAndDecode()
    if c8.debug {
      c8.debugInstruction(&instruction)
    } else {
      c8.executeInstruction(&instruction)
    }
    // TODO: run some tests to make sure this works like its supposed to, not sure
    // if the timing is right
    // TODO: also think if i want to edit the timing of key presses and how that works
    // maybe await key instruction i want on key press and release?
    if c8.delayTimer > 0 && loopCounter % int(c8.clockSpeed/DELAY_SOUND_TIMER_UPDATE) == 0 {
      c8.delayTimer -= 1
    }
    // TODO: implement sound!
    if c8.soundTimer > 0 && loopCounter % int(c8.clockSpeed/DELAY_SOUND_TIMER_UPDATE) == 0 {
      c8.soundTimer -= 1
    }
    time.Sleep(time.Duration(1000/c8.clockSpeed) * time.Millisecond)
    loopCounter += 1
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

  instructionMap := map[uint8]func(*utils.Instruction){}

  var debugState DebugState
  if debug {
    debugState = PAUSED
  } else {
    debugState = RUNNING
  }

  debugBreakpoint := uint16(0x0000)

  keyboard := new([16]keypress)

  c8 := Chip8{PROGRAM_START, 0, 0, 0, [32*64]uint8{}, utils.Stack{}, [16]uint8{}, memory, instructionMap, keyboard, CLOCK_SPEED, modern, debug, debugState, debugBreakpoint}

  // put instructions in a map
  c8.instructionMap[0x0] = c8.I0
  c8.instructionMap[0x1] = c8.I1NNN
  c8.instructionMap[0x2] = c8.I2NNN
  c8.instructionMap[0x3] = c8.I3XNN
  c8.instructionMap[0x4] = c8.I4XNN
  c8.instructionMap[0x5] = c8.I5XY0
  c8.instructionMap[0x6] = c8.I6XNN
  c8.instructionMap[0x7] = c8.I7XNN
  c8.instructionMap[0x8] = c8.I8XYN
  c8.instructionMap[0x9] = c8.I9XY0
  c8.instructionMap[0xA] = c8.IANNN
  c8.instructionMap[0xB] = c8.IBNNN
  c8.instructionMap[0xC] = c8.ICXNN
  c8.instructionMap[0xD] = c8.IDXYN
  c8.instructionMap[0xE] = c8.IE
  c8.instructionMap[0xF] = c8.IF

  return &c8
}
