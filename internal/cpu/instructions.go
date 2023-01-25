package cpu

import (
  "log"
  "math/rand"

  "jfeintzeig/chip8/internal/utils"
)

func (c8 *Chip8) I0(inst *utils.Instruction) {
  switch inst.NN {
  // 00E0: clear screen
  case 0xE0:
    c8.Display = [len(c8.Display)]uint8{}
  // 00EE: return from subroutine
  // pop stack and set pc to value
  case 0xEE:
    c8.stack, c8.pc = c8.stack.Pop()
  default:
    log.Fatal("unknown last nibble for I0 instruction.")
  }
}

// jump
func (c8 *Chip8) I1NNN(inst *utils.Instruction) {
  c8.pc = inst.NNN
}

// enter subroutine: store pc in stack and jump
func (c8 *Chip8) I2NNN(inst *utils.Instruction) {
  c8.stack = c8.stack.Push(c8.pc)
  c8.pc = inst.NNN
}

// skip instruction if VX == NN
func (c8 *Chip8) I3XNN(inst *utils.Instruction) {
  if c8.variableRegister[inst.X] == inst.NN {
    c8.pc += 2
  }
}

// skip instruction if VX != NN
func (c8 *Chip8) I4XNN(inst *utils.Instruction) {
  if c8.variableRegister[inst.X] != inst.NN {
    c8.pc += 2
  }
}

// skip instruction if VX == VY
func (c8 *Chip8) I5XY0(inst *utils.Instruction) {
  if c8.variableRegister[inst.X] == c8.variableRegister[inst.Y] {
    c8.pc += 2
  }
}

// skip instruction if VX != VY
func (c8 *Chip8) I9XY0(inst *utils.Instruction) {
  if c8.variableRegister[inst.X] != c8.variableRegister[inst.Y] {
    c8.pc += 2
  }
}

// set register VX to NN
func (c8 *Chip8) I6XNN(inst *utils.Instruction) {
  c8.variableRegister[inst.X] = inst.NN
}

// add NN to register VX
func (c8 *Chip8) I7XNN(inst *utils.Instruction) {
  c8.variableRegister[inst.X] += inst.NN
}

// logic and arithmetic
func (c8 *Chip8) I8XYN(inst *utils.Instruction) {
  switch inst.N {
  // 8XY0: set VX to VY
  case 0:
    c8.variableRegister[inst.X] = c8.variableRegister[inst.Y]
  // 8XY1: set VX to (VX | VY)
  case 1:
    c8.variableRegister[inst.X] |= c8.variableRegister[inst.Y]
    if !c8.modern {
      c8.variableRegister[0xF] = 0
    }
  // 8XY2: set VX to (VX & VY)
  case 2:
    c8.variableRegister[inst.X] &= c8.variableRegister[inst.Y]
    if !c8.modern {
      c8.variableRegister[0xF] = 0
    }
  // 8XY3: set VX to (VX XOR VY)
  case 3:
    c8.variableRegister[inst.X] ^= c8.variableRegister[inst.Y]
    if !c8.modern {
      c8.variableRegister[0xF] = 0
    }
  // 8XY4: set VX to (VX + VY)
  case 4:
    carry := uint8(0)
    if uint16(c8.variableRegister[inst.X]) + uint16(c8.variableRegister[inst.Y]) > 255 {
      carry = 1
    }
    c8.variableRegister[inst.X] += c8.variableRegister[inst.Y]
    c8.variableRegister[0xF] = carry
  // 8XY5: set VX to (VX - VY)
  case 5:
    carry := uint8(0)
    if c8.variableRegister[inst.X] > c8.variableRegister[inst.Y] {
      carry = 1
    }
    c8.variableRegister[inst.X] -= c8.variableRegister[inst.Y]
    c8.variableRegister[0xF] = carry
  // 8XY6: shift VX one bit right
  case 6:
    if !c8.modern {
      c8.variableRegister[inst.X] = c8.variableRegister[inst.Y]
    }

    rightMostBit := c8.variableRegister[inst.X] & 1
    c8.variableRegister[inst.X] = c8.variableRegister[inst.X] >> 1
    c8.variableRegister[0xF] = rightMostBit
  // 8XY7: set VX to (VY - VX)
  case 7:
    carry := uint8(0)
    if c8.variableRegister[inst.Y] > c8.variableRegister[inst.X] {
      carry = 1
    }
    c8.variableRegister[inst.X] = c8.variableRegister[inst.Y] - c8.variableRegister[inst.X]
    c8.variableRegister[0xF] = carry
  // 8XYE: shift VX one bit left
  case 0xE:
    if !c8.modern {
      c8.variableRegister[inst.X] = c8.variableRegister[inst.Y]
    }

    leftMostBit := (c8.variableRegister[inst.X] & 0x80) >> 7 // 0x80 is 10000000
    c8.variableRegister[inst.X] = c8.variableRegister[inst.X] << 1
    c8.variableRegister[0xF] = leftMostBit
  default:
    log.Fatal("unknown last nibble for I8XYN instruction.")
  }
}

// set index register to NNN
func (c8 *Chip8) IANNN(inst *utils.Instruction) {
  c8.i = inst.NNN
}

// jump to NNN + V0
func (c8 *Chip8) IBNNN(inst *utils.Instruction) {
  // TODO: the blog suggests non-modern is the preferred mode for this one,
  // but modern is the preferred mode for I8XYE and I8XY6? how to deal with this?
  if c8.modern {
   c8.pc = inst.NNN + uint16(c8.variableRegister[inst.X])
  } else {
   c8.pc = inst.NNN + uint16(c8.variableRegister[0])
  }
}

// random number, and with NN, put at VX
func (c8 *Chip8) ICXNN(inst *utils.Instruction) {
  c8.variableRegister[inst.X] = uint8(rand.Intn(256)) & inst.NN
}

// draw Display
func (c8 *Chip8) IDXYN(inst *utils.Instruction) {
  // choosing x+y starting place wraps the screen
  x := int(c8.variableRegister[inst.X] % 64)
  y := int(c8.variableRegister[inst.Y] % 32)
  c8.variableRegister[0xF] = 0
  for i := 0; i < int(inst.N); i++ {
    sprite := c8.memory[c8.i + uint16(i)]
    // for bit in sprite, loop over display and xor sprite bit and memory bit
    for spriteBit := 0; spriteBit < 8; spriteBit++ {
      // if part of the sprite is over the edge of the screen, clip it
      if ((x + spriteBit) < 64) && ((y + i) < 32) {
        // we have a 1d array representing a 2d screen; each 64 values is a row.
        index := (y+i)*64 + (x+spriteBit)
        // change display by xor'ing pixel with corresponding bit in sprite
        displayPixel := c8.Display[index]
        spritePixel := ((sprite >> (7-spriteBit)) & 0x01)
        if (displayPixel & spritePixel) == 1 {
          c8.variableRegister[0xF] = 1
        }
        c8.Display[index] = displayPixel ^ spritePixel
      }
    }
  }
}

// key presses
func (c8 *Chip8) IE(inst *utils.Instruction) {
  key := c8.variableRegister[inst.X]
  if key > 0xF {
    log.Fatalf("Unknown key %x", key)
  }

  switch inst.NN {
  // EX9E: if key corresponding to VX is pressed, skip next instruction
  case 0x9E:
    if c8.Keyboard[key].Pressed {
      c8.pc += 2
    }
  // EXA1: if key corresponding to VX is _not_pressed, skip next instruction
  case 0xA1:
    if !c8.Keyboard[key].Pressed {
      c8.pc += 2
    }
  default:
    log.Fatal("unknown last nibble for IE instruction.")
  }
}

// timers, fonts, keys, other stuff
func (c8 *Chip8) IF(inst *utils.Instruction) {
  switch inst.NN {
  // FX07: set VX to delay timer
  case 0x07:
    c8.variableRegister[inst.X] = c8.delayTimer
  // FX15: set delay timer to VX
  case 0x15:
    c8.delayTimer = c8.variableRegister[inst.X]
  // FX18: set sound timer to VX
  case 0x18:
    c8.soundTimer = c8.variableRegister[inst.X]
  // FX1E: add VX to index register
  case 0x1E:
    if uint16(c8.i) + uint16(c8.variableRegister[inst.X]) > 4095 {
      c8.variableRegister[0xF] = 1
    }
    c8.i += uint16(c8.variableRegister[inst.X])
  // FX0A: block until key press, then store key in VX
  case 0x0A:
    isKeyPressed := false
    key := uint8(0)
    for index, element := range c8.Keyboard {
      if element.Pressed {
        isKeyPressed = true
        key = uint8(index)
      }
    }
    if isKeyPressed {
    // Note: this blocks any cycles while it waits for key release
    // vs. the "wait for keypress" functionality of this function
    // actually returns so cycles run. The result is that timers
    // will continue decrementing until key press, then they will
    // pause until key release.
    OuterLoop:
      for true {
        if c8.Keyboard[key].JustReleased {
          break OuterLoop
        }
      }
    c8.variableRegister[inst.X] = key
    } else {
      c8.pc -= 2
    }
  // FX29: set index register to location of font character corresponding to last nibble of VX
  case 0x29:
    c8.i = FONT_START + 5 * uint16(c8.variableRegister[inst.X] & 0x0F)
  // FX33: take VX in decimal, separate each digit, and put in memory starting at index register
  // e.g. if VX is 0xAF, thats 175, so do:
  //     memory[i] = 1
  //     memory[i+1] = 7
  //     memory[i+2] = 5
  case 0x33:
    vx := c8.variableRegister[inst.X]
    c8.memory[c8.i] = vx / 100
    c8.memory[c8.i+1] = vx / 10 - (vx / 100)*10
    c8.memory[c8.i+2] = vx - (vx / 100)*100 - (vx / 10 - (vx / 100)*10)*10
  // FX55: Write variable register from V0 to VX (inclusive) into consecutive memory bytes, starting at index register address.
  case 0x55:
    for index := uint16(0); index <= uint16(inst.X); index++ {
      c8.memory[c8.i + index] = c8.variableRegister[index]
    }
    if !c8.modern {
      c8.i += uint16(inst.X) + 1
    }
  // FX65: Write X+1 consecutive memory bytes, starting at index register address, into variable register from V0 to VX, inclusive.
  case 0x65:
    for index := uint16(0); index <= uint16(inst.X); index++ {
      c8.variableRegister[index] = c8.memory[c8.i + index]
    }
    if !c8.modern {
      c8.i += uint16(inst.X) + 1
    }
  default:
    log.Fatal("uknown instruction %x for FX", inst.Full)
  }
}
