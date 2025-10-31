package ui

import (
	"synth/assets"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
)

type Ui struct {
	background *ebiten.Image
	w, h       int
	controls   *Controls
	menu       *Menu
}

func NewUi(asts *assets.Loader, ctrl *Controls, menu *preset.Node) (*Ui, error) {
	// BG + window size accordingly
	bg, err := asts.GetImage("ui/background")
	if err != nil {
		return nil, err
	}
	bds := bg.Bounds()
	ebiten.SetWindowSize(bds.Dx(), bds.Dy())

	// Main menu
	mu, err := NewMenu(asts, menu)
	if err != nil {
		return nil, err
	}

	ui := &Ui{
		background: bg,
		w:          bds.Dx(),
		h:          bds.Dy(),
		menu:       mu,
		controls:   ctrl,
	}

	return ui, nil
}

func (u *Ui) Update() error {
	_, _, s := u.controls.Update()
	u.menu.Scroll(s)
	u.menu.Update()
	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {
	screen.DrawImage(u.background, nil)
	u.menu.Draw(screen)
}

func (u *Ui) Layout(_, _ int) (int, int) {
	return u.w, u.h
}
