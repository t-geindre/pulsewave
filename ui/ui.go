package ui

import (
	"synth/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Ui struct {
	background     *ebiten.Image
	w, h           int
	menu           []*Entry
	selected       *ebiten.Image
	selectedXShift float64
	selectedYShift float64
	current        int
	controls       *Controls
}

func NewUi(asts *assets.Loader, ctrl *Controls) (*Ui, error) {
	bg, err := asts.GetImage("ui/background")
	if err != nil {
		return nil, err
	}

	bds := bg.Bounds()

	ebiten.SetWindowSize(bds.Dx(), bds.Dy())

	entries := make([]*Entry, 0)
	for _, name := range []string{"Oscillators", "Effects", "Filters", "Envelope", "Settings"} {
		entry, err := NewEntry(asts, name)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	selected, err := asts.GetImage("ui/selected")
	if err != nil {
		return nil, err
	}

	ui := &Ui{
		background: bg,
		w:          bds.Dx(),
		h:          bds.Dy(),
		menu:       entries,
		selected:   selected,
		current:    0,
		controls:   ctrl,
	}

	sbds := ui.selected.Bounds()
	ebds := ui.menu[0].Bounds()

	ui.selectedXShift = float64(sbds.Dx()-ebds.Dx()) / 2
	ui.selectedYShift = float64(sbds.Dy()-ebds.Dy()) / 2

	return ui, nil
}

func (u *Ui) Update() error {
	u.controls.Update()
	_, _, m, l := u.controls.Consume()
	u.current -= m - l
	u.current = (u.current + len(u.menu)) % len(u.menu)
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
