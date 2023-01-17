package display

import (
  //"fmt"
  "image/color"
  //"time"
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
  debug *bool
}

func (g *Game) Update() error {
  //time.Sleep(time.Duration(1000/g.c8.ClockSpeed) * time.Millisecond)
  instruction := g.c8.FetchAndDecode()
  g.c8.Execute(&instruction)
  //if *g.debug {
  //  fmt.Println(g.c8)
    // fmt.Println("Press Enter to continue")
    // bufio.NewReader(os.Stdin).ReadBytes('\n')
  //}
  return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  //ebitenutil.DebugPrint(screen, "Hello, World!")
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

func NewGame(c8 *cpu.Chip8, debug *bool) (*Game, error) {
	g := &Game{
    c8,
    debug,
	}
	return g, nil
}
