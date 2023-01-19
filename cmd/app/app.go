package main

import (
  "flag"
  "fmt"
  "log"

  "github.com/hajimehoshi/ebiten/v2"

  "jfeintzeig/chip8/internal/cpu"
  "jfeintzeig/chip8/internal/display"
)

var (
  debug *bool
  modern *bool
  file *string
)

func init() {
  file = flag.String("file","data/ibm_logo.ch8","path to file to load")
  debug = flag.Bool("debug",false,"set true to debug output")
  modern = flag.Bool("modern",true,"set true to use modern interpretation of ambiguous instructions, default true")
}

func main() {
  flag.Parse()

  fmt.Println("Starting up...")
  chip8 := cpu.NewChip8(*debug, *modern)
  chip8.LoadFile(*file)

  ebiten.SetWindowSize(640, 320)
  ebiten.SetWindowTitle("Hello, World!")
  game, _ := display.NewGame(chip8)

  // infinite loop at chip8.clockSpeed
  go chip8.Execute()

  // display updates @ 60Hz via infinite loop in ebiten
  if err := ebiten.RunGame(game); err != nil {
    log.Fatal(err)
  }
}
