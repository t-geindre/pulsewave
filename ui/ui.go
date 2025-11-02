package ui

import (
	"fmt"
	"synth/assets"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
)

// Todo get it fron config
const (
	BodyHeight = 211
	BodyWidth  = 376
	BodyStartX = 51
	BodyStartY = 70
)

type Ui struct {
	background   *ebiten.Image
	w, h         int
	controls     *Controls
	component    Component
	bodyClipMask *ebiten.Image
}

func NewUi(asts *assets.Loader, ctrl *Controls, menu *preset.Node) (*Ui, error) {
	// BG + window size accordingly
	bg, err := asts.GetImage("ui/background")
	if err != nil {
		return nil, err
	}
	bds := bg.Bounds()
	ebiten.SetWindowSize(bds.Dx(), bds.Dy())

	// Main component
	mu, err := NewList(asts, menu)
	if err != nil {
		return nil, err
	}

	ui := &Ui{
		background:   bg,
		w:            bds.Dx(),
		h:            bds.Dy(),
		component:    mu,
		controls:     ctrl,
		bodyClipMask: ebiten.NewImage(BodyWidth, BodyHeight),
	}

	return ui, nil
}

func (u *Ui) Update() error {
	fw, _, s := u.controls.Update()
	if fw {
		fmt.Println(u.component.CurrentTarget().Label)
	}
	u.component.Scroll(s)
	u.component.Update()
	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {
	screen.DrawImage(u.background, nil)

	u.bodyClipMask.Clear()
	u.component.Draw(u.bodyClipMask)

	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(BodyStartX, BodyStartY)
	screen.DrawImage(u.bodyClipMask, ops)
}

func (u *Ui) Layout(_, _ int) (int, int) {
	return u.w, u.h
}
