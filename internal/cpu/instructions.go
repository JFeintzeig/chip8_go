package cpu

import "fmt"

// clear screen
func (c8 *Chip8) I00E0(inst *instruction) {
  if inst.full == 0x00e0 {
    c8.Display = [len(c8.Display)]uint8{}
  }
}

// jump
func (c8 *Chip8) I1NNN(inst *instruction) {
  c8.pc = inst.nnn
}

// set register VX to NN
func (c8 *Chip8) I6XNN(inst *instruction) {
  c8.variableRegister[inst.x] = inst.nn
}

// add NN to register VX
func (c8 *Chip8) I7XNN(inst *instruction) {
  c8.variableRegister[inst.x] += inst.nn
}

// set index register to NNN
func (c8 *Chip8) IANNN(inst *instruction) {
  c8.i = inst.nnn
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
        c8.Display[index] = c8.Display[index] ^ (sprite >> (7-spriteBit) & 0x01)
      }
    }
    // TODO: refactor this somewhere else
    if c8.debug {
      fmt.Printf(" x %d \n y %d \n display index %d \n sprite %x \n",x,y,(y+i)*64 + x,sprite)
      fmt.Println(c8.Display)
    }
  }
}
