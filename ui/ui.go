package ui

import (
	"synth/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Ui struct {
	background     *ebiten.Image
	w, h           int
	menu           []*Entry
	selected       *ebiten.Image
	selectedXShift float64
	selectedYShift float64
	current        int
}

func NewUi(asts *assets.Loader) *Ui {
	bg := asts.GetImage("ui/background")
	bds := bg.Bounds()

	ebiten.SetWindowSize(bds.Dx(), bds.Dy())

	ui := &Ui{
		background: bg,
		w:          bds.Dx(),
		h:          bds.Dy(),
		menu: []*Entry{
			NewEntry(asts, "Oscillators"),
			NewEntry(asts, "Effects"),
			NewEntry(asts, "Filters"),
			NewEntry(asts, "Envelope"),
			NewEntry(asts, "Settings"),
		},
		selected: asts.GetImage("ui/selected"),
		current:  0,
	}

	sbds := ui.selected.Bounds()
	ebds := ui.menu[0].Bounds()

	ui.selectedXShift = float64(sbds.Dx()-ebds.Dx()) / 2
	ui.selectedYShift = float64(sbds.Dy()-ebds.Dy()) / 2

	return ui
}

func (u *Ui) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		u.current--
		if u.current < 0 {
			u.current = len(u.menu) - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		u.current++
		if u.current >= len(u.menu) {
			u.current = 0
		}
	}
	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {
	const Items = 4 // Todo get it from config
	const MenuStartX = 76
	const MenuStartY = 42
	const EntrySpacing = 50

	screen.DrawImage(u.background, nil)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(MenuStartX, MenuStartY)

	for i := u.current - 1; i < u.current+Items-1; i++ {
		p := (i + len(u.menu)) % len(u.menu)
		entry := u.menu[p]

		opts.GeoM.Translate(0, EntrySpacing)
		screen.DrawImage(entry.Image, opts)

		if i == u.current {
			// Draw selected background
			selOpts := &ebiten.DrawImageOptions{}
			selOpts.GeoM.Concat(opts.GeoM)
			selOpts.GeoM.Translate(-u.selectedXShift, -u.selectedYShift)
			screen.DrawImage(u.selected, selOpts)
		}
	}
}

func (u *Ui) Layout(_, _ int) (int, int) {
	return u.w, u.h
}

func (u *Ui) loadResources() error {
	return nil
}
