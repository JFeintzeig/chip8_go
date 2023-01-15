package cpu

// clear screen
func (c8 *chip8) I00E0(inst *instruction) {
  if inst.full == 0x00e0 {
    c8.display = [len(c8.display)]uint8{}
  }
}

// jump
func (c8 *chip8) I1NNN(inst *instruction) {
  c8.pc = inst.nnn
}

// set register VX to NN
func (c8 *chip8) I6XNN(inst *instruction) {
  c8.variableRegister[inst.x] = inst.nn
}

// add NN to register VX
func (c8 *chip8) I7XNN(inst *instruction) {
  c8.variableRegister[inst.x] += inst.nn
}

// set index register to NNN
func (c8 *chip8) IANNN(inst *instruction) {
  c8.i = inst.nnn
}

// draw display
func (c8 *chip8) IDXYN(inst * instruction) {
}
