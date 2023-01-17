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
  file *string
)

func init() {
  file = flag.String("file","data/ibm_logo.ch8","path to file to load")
  debug = flag.Bool("debug",false,"set true to debug output")
}

func main() {
  flag.Parse()

  fmt.Println("Starting up...")
  chip8 := cpu.NewChip8(*debug)
  chip8.LoadFile(*file)

  // debug
  if *debug {
    chip8.SetDisplay()
    fmt.Println(chip8)
  }

  ebiten.SetWindowSize(640, 320)
  ebiten.SetWindowTitle("Hello, World!")
  game, _ := display.NewGame(chip8, debug)
  // do i want this? do i want to configure loop differently?
  ebiten.SetTPS(500)
  // execution loop now runs in ebiten
  if err := ebiten.RunGame(game); err != nil {
    log.Fatal(err)
  }
}
