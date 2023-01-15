package main

import (
  "fmt"
  "log"
  "os"
  "time"
  "jfeintzeig/chip8/internal/cpu"
)

var debug bool = true

func main() {
  fmt.Println("Starting up...")
  file, err := os.Open("data/ibm_logo.ch8")
  if err != nil {
    log.Fatal("can't find file")
  }
  chip8 := cpu.NewChip8()
  // debug
  if debug {
    chip8.SetDisplay()
    fmt.Println(chip8)
  }
  for true {
    time.Sleep(1 * time.Second)
    instruction := chip8.FetchAndDecode(*file)
    chip8.Execute(&instruction)
    if debug {
      fmt.Println(chip8)
    }
  }
}
