package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	width  = 1280
	height = 720
)

var (
	deepField *ebiten.Image
	isWasm    = runtime.GOARCH == "wasm"
)

// Vieport code taken from: https://github.com/hajimehoshi/ebiten/blob/master/examples/infinitescroll/main.go
// Credit goes to The Ebiten Authors
type viewport struct {
	x16 int
	y16 int
}

func (p *viewport) Move() {
	w, h := deepField.Size()
	maxX16 := w * 16
	maxY16 := h * 16

	p.x16 += w / 32
	p.y16 += h / 32
	p.x16 %= maxX16
	p.y16 %= maxY16
}

func (p *viewport) Position() (int, int) {
	return p.x16, p.y16
}

type Game struct {
	viewport viewport
}

func (g *Game) Update() error {
	if !isWasm && ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("esc pressed")
	}

	g.viewport.Move()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x16, y16 := g.viewport.Position()
	offsetX, offsetY := float64(-x16)/16, float64(-y16)/16

	const repeat = 3
	w, h := deepField.Size()
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(w*i), float64(h*j))
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(deepField, op)
		}
	}

	msg := fmt.Sprintf(`TPS: %0.2f; FPS: %0.2f
vp: %d, %d
`, ebiten.CurrentTPS(), ebiten.CurrentFPS(),
		g.viewport.x16, g.viewport.y16)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(Deepfield_png))
	if err != nil {
		log.Fatal(err)
	}

	deepField = ebiten.NewImageFromImage(img)
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Deepfield scroll")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
