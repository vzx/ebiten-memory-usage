package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/vzx/ebiten-memory-usage/assets"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	layoutWidth, layoutHeight = 3840, 2160
	windowWidth, windowHeight = 1920, 1080
)

var (
	deepField *ebiten.Image
	mplusBold font.Face
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
	start        time.Time
	totalRunTime time.Duration
	viewport     viewport
	ticks        uint64
	memStats     *runtime.MemStats
}

func (g *Game) Update() error {
	if !isWasm && ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("esc pressed")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		runtime.GC()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		max := ebiten.MaxTPS()
		if max == ebiten.UncappedTPS {
			ebiten.SetMaxTPS(60)
		} else {
			ebiten.SetMaxTPS(ebiten.UncappedTPS)
		}
	}

	g.ticks++
	if g.ticks%30 == 0 {
		runtime.ReadMemStats(g.memStats)
		g.totalRunTime = time.Now().Sub(g.start)
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

	ms := g.memStats
	msg := fmt.Sprintf(`vp: %d, %d
TPS: %0.2f (max: %d); FPS: %0.2f
Run time: %v
ticks: %d
Alloc: %s
Total: %s
Sys: %s
NextGC: %s
NumGC: %d

<G>: run garbage collection
<T>: toggle unlimited TPS`,
		g.viewport.x16, g.viewport.y16,
		ebiten.CurrentTPS(), ebiten.MaxTPS(), ebiten.CurrentFPS(),
		g.totalRunTime,
		g.ticks,
		formatBytes(ms.Alloc), formatBytes(ms.TotalAlloc), formatBytes(ms.Sys),
		formatBytes(ms.NextGC), ms.NumGC,
	)
	text.Draw(screen, msg, mplusBold, 11, 55, color.Black)
	text.Draw(screen, msg, mplusBold, 10, 54, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return layoutWidth, layoutHeight
}

func formatBytes(b uint64) string {
	if b >= 1073741824 {
		return fmt.Sprintf("%0.2f GiB", float64(b)/1073741824)
	} else if b >= 1048576 {
		return fmt.Sprintf("%0.2f MiB", float64(b)/1048576)
	} else if b >= 1024 {
		return fmt.Sprintf("%0.2f KiB", float64(b)/1024)
	} else {
		return fmt.Sprintf("%d B", b)
	}
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(assets.Deepfield_png))
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(assets.Mplus_1p_bold_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusBold, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	deepField = ebiten.NewImageFromImage(img)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Deepfield scroll")

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	if err := ebiten.RunGame(&Game{memStats: memStats, start: time.Now()}); err != nil {
		log.Fatal(err)
	}
}
