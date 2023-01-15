package main

import (
  "fmt"
  "log"
  "os"
  "time"
  "jfeintzeig/chip8/internal/cpu"
)

var debug bool = false

func main() {
  fmt.Println("Starting up...")
  file, err := os.Open("data/ibm_logo.ch8")
  if err != nil {
    log.Fatal("can't find file")
  }
  fmt.Println(cpu.PC)
  fmt.Println(cpu.I)
  fmt.Println(cpu.Font)
  fmt.Println(cpu.Memory)
  fmt.Println(*file)
  for true {
    time.Sleep(1 * time.Second)
    instruction := cpu.Fetch(*file)
    cpu.Execute(instruction)
    if debug {
      fmt.Println(cpu.Display)
    }
  }
}
