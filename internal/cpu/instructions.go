package cpu

import (
  "log"
)

// clear screen
func (c8 *Chip8) I00E0(inst *instruction) {
  if inst.full == 0x00e0 {
    c8.Display = [len(c8.Display)]uint8{}
  }
}

// TODO: deal with decode/route for dupes like 0
func (c8 *Chip8) I0NNN(inst *instruction) {
  log.Fatal("I0NNN is not implemented")
}

// TODO: implement
func (c8 *Chip8) I00EE(inst *instruction) {
}

// jump
func (c8 *Chip8) I1NNN(inst *instruction) {
  c8.pc = inst.nnn
}

// TODO: implement
func (c8 *Chip8) I2NNN(inst *instruction) {
}

// skip instruction if VX == NN
func (c8 *Chip8) I3XNN(inst *instruction) {
  if c8.variableRegister[inst.x] == inst.nn {
    c8.pc += 2
  }
}

// skip instruction if VX != NN
func (c8 *Chip8) I4XNN(inst *instruction) {
  if c8.variableRegister[inst.x] != inst.nn {
    c8.pc += 2
  }
}

// skip instruction if VX == VY
func (c8 *Chip8) I5XY0(inst *instruction) {
  if c8.variableRegister[inst.x] == c8.variableRegister[inst.y] {
    c8.pc += 2
  }
}

// skip instruction if VX != VY
func (c8 *Chip8) I9XY0(inst *instruction) {
  if c8.variableRegister[inst.x] != c8.variableRegister[inst.y] {
    c8.pc += 2
  }
}

// set register VX to NN
func (c8 *Chip8) I6XNN(inst *instruction) {
  c8.variableRegister[inst.x] = inst.nn
}

// add NN to register VX
func (c8 *Chip8) I7XNN(inst *instruction) {
  c8.variableRegister[inst.x] += inst.nn
}

// logic and arithmetic
func (c8 *Chip8) I8XYN(inst *instruction) {
  switch inst.n {
  // set VX to VY
  case 0:
    c8.variableRegister[inst.x] = c8.variableRegister[inst.y]
  // set VX to (VX | VY)
  case 1:
    c8.variableRegister[inst.x] |= c8.variableRegister[inst.y]
  // set VX to (VX & VY)
  case 2:
    c8.variableRegister[inst.x] &= c8.variableRegister[inst.y]
  // set VX to (VX XOR VY)
  case 3:
    c8.variableRegister[inst.x] ^= c8.variableRegister[inst.y]
  // set VX to (VX + VY)
  case 4:
    if uint16(c8.variableRegister[inst.x]) + uint16(c8.variableRegister[inst.y]) > 255 {
      c8.variableRegister[0xF] = 1
    } else {
      c8.variableRegister[0xF] = 0
    }
    c8.variableRegister[inst.x] += c8.variableRegister[inst.y]
  // set VX to (VX - VY)
  case 5:
    if c8.variableRegister[inst.x] > c8.variableRegister[inst.y] {
      c8.variableRegister[0xF] = 1
    } else {
      c8.variableRegister[0xF] = 0
    }
    c8.variableRegister[inst.x] -= c8.variableRegister[inst.y]
  // shift VX one bit right
  case 6:
    if !c8.modern {
      c8.variableRegister[inst.x] = c8.variableRegister[inst.y]
    }

    rightMostBit := c8.variableRegister[inst.x] & 1
    c8.variableRegister[inst.x] = c8.variableRegister[inst.x] >> 1
    c8.variableRegister[0xF] = rightMostBit
  // set VX to (VY - VX)
  case 7:
    if c8.variableRegister[inst.y] > c8.variableRegister[inst.x] {
      c8.variableRegister[0xF] = 1
    } else {
      c8.variableRegister[0xF] = 0
    }
    c8.variableRegister[inst.x] = c8.variableRegister[inst.y] - c8.variableRegister[inst.x]
  // shift VX one bit left
  case 0xE:
    if !c8.modern {
      c8.variableRegister[inst.x] = c8.variableRegister[inst.y]
    }

    leftMostBit := (c8.variableRegister[inst.x] & 0x80) >> 7 // 0x80 is 10000000
    c8.variableRegister[inst.x] = c8.variableRegister[inst.x] << 1
    c8.variableRegister[0xF] = leftMostBit
  default:
    log.Fatal("unknown last nibble for I8XYN instruction.")
  }
}

// set index register to NNN
func (c8 *Chip8) IANNN(inst *instruction) {
  c8.i = inst.nnn
}

// jump to NNN + V0
func (c8 *Chip8) IBNNN(inst *instruction) {
  // TODO: the blog suggests non-modern is the preferred mode for this one,
  // but modern is the preferred mode for I8XYE and I8XY6? how to deal with this?
  if c8.modern {
   c8.pc = inst.nnn + uint16(c8.variableRegister[inst.x])
  } else {
   c8.pc = inst.nnn + uint16(c8.variableRegister[0])
  }
}

// draw Display
func (c8 *Chip8) IDXYN(inst * instruction) {
  // choosing x+y starting place wraps the screen
  x := int(c8.variableRegister[inst.x] % 64)
  y := int(c8.variableRegister[inst.y] % 32)
  // do we draw down or up?
  for i := 0; i < int(inst.n); i++ {
    sprite := c8.memory[c8.i + uint16(i)]
    // for bit in sprite, loop over display and xor sprite bit and memory bit
    for spriteBit := 0; spriteBit < 8; spriteBit++ {
      // if part of the sprite is over the edge of the screen, clip it
      if ((x + spriteBit) < 64) && ((y + i) < 32) {
        // we have a 1d array representing a 2d screen; each 64 values is a row.
        index := (y+i)*64 + (x+spriteBit)
        // change display by xor'ing pixel with corresponding bit in sprite
        displayPixel := c8.Display[index]
        spritePixel := (sprite >> (7-spriteBit) & 0x01)
        if (displayPixel & spritePixel) == 1 {
          c8.variableRegister[0xF] = 1
        } else {
          c8.variableRegister[0xF] = 0
        }
        c8.Display[index] = displayPixel ^ spritePixel
      }
    }
  }
}
