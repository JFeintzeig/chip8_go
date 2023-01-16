package display

import (
  "fmt"
  "image/color"
  "time"
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/ebitenutil"
  "jfeintzeig/chip8/internal/cpu"
)

var (
  pixel = ebiten.NewImage(10,10)
)

func init() {
  pixel.Fill(color.White)
}

type Game struct {
  c8 *cpu.Chip8
  debug *bool
}

func (g *Game) Update() error {
  // TODO: how do i control speed? looks like ebiten built in has different
  // speeds for Update() and Draw(), i think i want Update() at ClockSpeed
  // and Draw() at 60Hz?
  time.Sleep(time.Duration(1000/g.c8.ClockSpeed) * time.Millisecond)
  instruction := g.c8.FetchAndDecode()
  g.c8.Execute(&instruction)
  if *g.debug {
    fmt.Println(g.c8)
    // fmt.Println("Press Enter to continue")
    // bufio.NewReader(os.Stdin).ReadBytes('\n')
  }
  return nil
}

// TODO: need a function to translate from g.c8.Display into screen
func (g *Game) Draw(screen *ebiten.Image) {
  ebitenutil.DebugPrint(screen, "Hello, World!")
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(40,40)
  screen.DrawImage(pixel, op)

  op = &ebiten.DrawImageOptions{}
  op.GeoM.Translate(40,80)
  screen.DrawImage(pixel, op)
  op = &ebiten.DrawImageOptions{}
  op.GeoM.Translate(50,80)
  screen.DrawImage(pixel, op)

  op = &ebiten.DrawImageOptions{}
  op.GeoM.Translate(40,51)
  screen.DrawImage(pixel, op)
  op = &ebiten.DrawImageOptions{}
  op.GeoM.Translate(42,51)
  screen.DrawImage(pixel, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
  // TODO: abstract this away.
  // rn layout and `pixel` are 10x actual memory in Chip8, need to codify this somewhere better...
  return 640, 320
}

func NewGame(c8 *cpu.Chip8, debug *bool) (*Game, error) {
	g := &Game{
    c8,
    debug,
	}
	return g, nil
}
