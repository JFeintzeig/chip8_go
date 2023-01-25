package main

import (
  "flag"
  "fmt"

  "jfeintzeig/chip8/internal/disassembler"
)

var (
  inputFile *string
  outputFile *string
)

func init() {
  inputFile = flag.String("inputFile","data/ibm_logo.ch8","path to file to disassemble")
  outputFile = flag.String("outputFile","out.8o","where to save output")
}

func main() {
  flag.Parse()

  fmt.Println("Disassembling...")

  dis := disassembler.NewDisassembler(inputFile, outputFile)
  dis.Disassemble()
}
