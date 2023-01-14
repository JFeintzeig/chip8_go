package main

import "fmt"
import "jfeintzeig/chip8/internal/utils"
import "jfeintzeig/chip8/internal/cpu"

func main() {
    fmt.Println(utils.GetText())
    fmt.Println(cpu.PC)
    fmt.Println(cpu.I)
    fmt.Println(cpu.Font)
}
