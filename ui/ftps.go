package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Ftps struct {
}

func NewFtps() *Ftps {
	return &Ftps{}
}

func (f *Ftps) Draw(screen *ebiten.Image) {
	const (
		padding = 5
		cw      = 6
		ch      = 16
	)

	str := fmt.Sprintf("FPS: %.0f TPS: %.0f", ebiten.ActualFPS(), ebiten.ActualTPS())
	length := float32(len(str) * cw)

	vector.DrawFilledRect(
		screen,
		0, 0,
		length+padding*2,
		ch+padding*2,
		color.RGBA{R: 0, G: 0, B: 0, A: 100},
		false,
	)

	ebitenutil.DebugPrintAt(screen, str, padding, padding)
}
