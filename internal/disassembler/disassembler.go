package disassembler

import (
  "bufio"
  "errors"
  "fmt"
  "io"
  "log"
  "os"

  "jfeintzeig/chip8/internal/utils"
)

type disassembler struct {
  inputFile *string
  outputFile *string
  instructionStrings map[uint8] func(*utils.Instruction) string
}

func (dis *disassembler) Disassemble() {
  inputFile, err := os.Open(*dis.inputFile)
  if err != nil {
    log.Fatal("can't find inputFile")
  }
  br := bufio.NewReader(inputFile)

  outputFile, err := os.Create(*dis.outputFile)
  defer outputFile.Close()

  for {
    buf := make([]byte, 2)
    _, err := io.ReadFull(br, buf)
    if err != nil && !errors.Is(err, io.EOF) {
        fmt.Println(err)
        break
    }
    if err != nil {
        // end of inputFile
        break
    }
    inst := utils.DecodeInstruction((uint16(buf[0]) << 8) & uint16(buf[1]))
    outputFile.WriteString(dis.instructionStrings[inst.A](&inst) + "\n")
  }
}

func (dis *disassembler) I0(inst *utils.Instruction) string {
  return "test"
}

func NewDisassembler(inputFile *string, outputFile *string) disassembler {

  instructionStrings := map[uint8] func(*utils.Instruction) string {}

  dis := disassembler{
    inputFile: inputFile,
    outputFile: outputFile,
    instructionStrings: instructionStrings,
  }

  dis.instructionStrings[0x0] = dis.I0

  return dis
}
