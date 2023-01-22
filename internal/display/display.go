package display

import (
  "image/color"
  "github.com/hajimehoshi/ebiten/v2/audio"
  "github.com/hajimehoshi/ebiten/v2/audio/mp3"
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/ebitenutil"
  "github.com/hajimehoshi/ebiten/v2/inpututil"
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
  keyboard [16]ebiten.Key
  audioPlayer *audio.Player
}

func (g *Game) Update() error {
  for index, key := range g.keyboard {
    g.c8.Keyboard[index].Pressed = ebiten.IsKeyPressed(key)
    g.c8.Keyboard[index].JustReleased = inpututil.IsKeyJustReleased(key)
  }
  if g.c8.GetSoundTimer() > 0 {
    g.audioPlayer.Play()
    g.audioPlayer.Rewind()
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
    ebiten.KeyX,
    ebiten.KeyDigit1,
    ebiten.KeyDigit2,
    ebiten.KeyDigit3,
    ebiten.KeyQ,
    ebiten.KeyW,
    ebiten.KeyE,
    ebiten.KeyA,
    ebiten.KeyS,
    ebiten.KeyD,
    ebiten.KeyZ,
    ebiten.KeyC,
    ebiten.KeyDigit4,
    ebiten.KeyR,
    ebiten.KeyF,
    ebiten.KeyV,
  }

  audioContext := audio.NewContext(48000)
  f, _ := ebitenutil.OpenFile("data/beep.mp3")
  d, _ := mp3.Decode(audioContext, f)
  audioPlayer, _ := audio.NewPlayer(audioContext, d)

  g := &Game{
    c8,
    keyboard,
    audioPlayer,
  }
  return g, nil
}
