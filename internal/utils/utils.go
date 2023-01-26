package utils

import (
  "log"
  "text/template"
)

// represents the free parameters in an instruction
type InstructionType int

const (
  FULL InstructionType = iota
  X
  XY
  XYN
  XNN
  NNN
)

type Instruction struct {
  Full uint16
  A uint8
  X uint8
  Y uint8
  N uint8
  NN uint8
  NNN uint16
  Mnemonic string
  Type InstructionType
  Template *template.Template
}

// TODO: maybe set templates for all these, and then do we need Type + Mnemonic?
// but think about how we'll want to do the reverse (e.g. for assembler, string -> bytecode)
func (inst *Instruction) SetMnemonicTypeTemplate() {
  switch inst.A {
  case 0x0:
    switch inst.Full {
    case 0x00E0:
      inst.Mnemonic = "clear"
      inst.Type = FULL
    case 0x00EE:
      inst.Mnemonic = "return"
      inst.Type = FULL
    }
  case 0x1:
    inst.Mnemonic = "jump"
    inst.Type = NNN
  case 0x2:
    inst.Mnemonic = "call"
    inst.Type = NNN
  case 0x3:
    inst.Mnemonic = "skipe"
    inst.Type = XNN
  case 0x4:
    inst.Mnemonic = "skipne"
    inst.Type = XNN
  case 0x5:
    inst.Mnemonic = "skipre"
    inst.Type = XY
  case 0x6:
    inst.Mnemonic = "set"
    inst.Type = XNN
  case 0x7:
    inst.Mnemonic = "add"
    inst.Type = XNN
  case 0x8:
    switch inst.N {
    case 0x0:
      inst.Mnemonic = "move"
      inst.Type = XY
    case 0x1:
      inst.Mnemonic = "or"
      inst.Type = XY
    case 0x2:
      inst.Mnemonic = "and"
      inst.Type = XY
    case 0x3:
      inst.Mnemonic = "xor"
      inst.Type = XY
    case 0x4:
      inst.Mnemonic = "addr"
      inst.Type = XY
    case 0x5:
      inst.Mnemonic = "sub"
      inst.Type = XY
    case 0x6:
      inst.Mnemonic = "shiftr"
      inst.Type = XY
    case 0x7:
      inst.Mnemonic = "subr"
      inst.Type = XY
    case 0xE:
      inst.Mnemonic = "shiftl"
      inst.Type = XY
    }
  case 0x9:
    inst.Mnemonic = "skiprne"
    inst.Type = XY
  case 0xA:
    inst.Mnemonic = "seti"
    inst.Type = NNN
  case 0xB:
    inst.Mnemonic = "jump0"
    inst.Type = NNN
  case 0xC:
    inst.Mnemonic = "rand"
    inst.Type = XNN
  case 0xD:
    inst.Mnemonic = "sprite"
    inst.Type = XYN
  case 0xE:
    switch inst.NN {
    case 0x9E:
      inst.Mnemonic = "skipkey"
      inst.Type = X
    case 0xA1:
      inst.Mnemonic = "skipnkey"
      inst.Type = X
    }
  case 0xF:
    switch inst.NN {
    case 0x07:
      inst.Mnemonic = "setfromdelay"
      inst.Type = X
    case 0x15:
      inst.Mnemonic = "settodelay"
      inst.Type = X
    case 0x18:
      inst.Mnemonic = "settosound"
      inst.Type = X
    case 0x1E:
      inst.Mnemonic = "addi"
      inst.Type = X
    case 0x0A:
      inst.Mnemonic = "key"
      inst.Type = X
    case 0x29:
      inst.Mnemonic = "font"
      inst.Type = X
    case 0x33:
      inst.Mnemonic = "bcd"
      inst.Type = X
    case 0x55:
      inst.Mnemonic = "save"
      inst.Type = X
    case 0x65:
      inst.Mnemonic = "save"
      inst.Type = X
    }
  }
}

func (inst *Instruction) ToString() string {
  // TODO: use StringTemplate, this should be the same for all instructions
  return "not implemented"
}

func (inst *Instruction) ToBytecode() uint16 {
  // TODO: switch statement over inst.Type
  return 0
}

type Stack []uint16

func (s Stack) Push(val uint16) Stack {
  if len(s) > 16 {
    log.Fatal("stack overflow")
  }
  return append(s, val)
}

func (s Stack) Pop() (Stack, uint16) {
  l := len(s)
  last := s[l-1]
  s = s[:l-1]
  return s, last
}

// given bytecode, decode and parse instruction
func InstructionFromBytecode(codedInstruction uint16) Instruction {
  // TODO: infer instruction type, mnemonic, template
  inst := Instruction{
    Full: codedInstruction,
    A: uint8((codedInstruction & 0xF000) >> 12),
    X: uint8((codedInstruction & 0x0F00) >> 8),
    Y: uint8((codedInstruction & 0x00F0) >> 4),
    N: uint8(codedInstruction & 0x000F),
    NN: uint8(codedInstruction & 0x00FF),
    NNN: codedInstruction & 0x0FFF,
  }
  inst.SetMnemonicTypeTemplate()

  switch inst.Type {
    case FULL:
      inst.Template, _ = template.New("Full").Parse("{{.Mnemonic}} # {{printf \"%04X\" .Full}}")
    case X:
      inst.Template, _ = template.New("XY").Parse("{{.Mnemonic}} V{{printf \"%X\" .X}} # {{printf \"%04X\" .Full}}")
    case XY:
      inst.Template, _ = template.New("XY").Parse("{{.Mnemonic}} V{{printf \"%X\" .X}} V{{printf \"%X\" .Y}} # {{printf \"%04X\" .Full}}")
    case XYN:
      inst.Template, _ = template.New("XYN").Parse("{{.Mnemonic}} V{{printf \"%X\" .X}} V{{printf \"%X\" .Y}} {{printf \"%X\" .N}} # {{printf \"%04X\" .Full}}")
    case XNN:
      inst.Template, _ = template.New("XNN").Parse("{{.Mnemonic}} V{{printf \"%X\" .X}} {{printf \"%02X\" .NN}} # {{printf \"%04X\" .Full}}")
    case NNN:
      inst.Template, _ = template.New("NNN").Parse("{{.Mnemonic}} {{printf \"%04X\" .NNN}} # {{printf \"%04X\" .Full}}")
  }

  return inst
}

// another method for InstructionFromString?
