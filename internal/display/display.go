package display

import (
  "image/color"
  "github.com/hajimehoshi/ebiten/v2"
  //"github.com/hajimehoshi/ebiten/v2/ebitenutil"
  "jfeintzeig/chip8/internal/cpu"
)

var (
  // TODO: how do i want it to look? maybe change to 10, could add gridlines for debugging.
  pixel = ebiten.NewImage(9,10)
)

func init() {
  pixel.Fill(color.White)
}

type Game struct {
  c8 *cpu.Chip8
  keyboard [16]ebiten.Key
}

func (g *Game) Update() error {
  for index, key := range g.keyboard {
    g.c8.Keyboard[index] = ebiten.IsKeyPressed(key)
  }
  return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  for index, element := range g.c8.Display {
    if element == 1 {
      op := &ebiten.DrawImageOptions{}
      y := int(index / 64)
      x := int(index % 64)
      op.GeoM.Translate(float64(x*10),float64(y*10))
      screen.DrawImage(pixel, op)
    }
  }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
  // TODO: abstract this away.
  // rn layout and `pixel` are 10x actual memory in Chip8, need to codify this somewhere better...
  return 640, 320
}

func NewGame(c8 *cpu.Chip8) (*Game, error) {
  keyboard := [16]ebiten.Key{
    ebiten.KeyDigit1,
    ebiten.KeyDigit2,
    ebiten.KeyDigit3,
    ebiten.KeyQ,
    ebiten.KeyW,
    ebiten.KeyE,
    ebiten.KeyA,
    ebiten.KeyS,
    ebiten.KeyD,
    ebiten.KeyX,
    ebiten.KeyZ,
    ebiten.KeyC,
    ebiten.KeyDigit4,
    ebiten.KeyR,
    ebiten.KeyF,
    ebiten.KeyV,
  }
	g := &Game{
    c8,
    keyboard,
	}
	return g, nil
}
