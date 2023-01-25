package utils

import (
  "log"
)

type Instruction struct {
  Full uint16
  A uint8
  X uint8
  Y uint8
  N uint8
  NN uint8
  NNN uint16
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

func DecodeInstruction(codedInstruction uint16) Instruction {
  return Instruction{
    codedInstruction,
    uint8((codedInstruction & 0xF000) >> 12),
    uint8((codedInstruction & 0x0F00) >> 8),
    uint8((codedInstruction & 0x00F0) >> 4),
    uint8(codedInstruction & 0x000F),
    uint8(codedInstruction & 0x00FF),
    codedInstruction & 0x0FFF,
  }
}
