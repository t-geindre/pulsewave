package ui

import (
	"errors"
	"synth/assets"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// Todo get it from config
	MenuItems        = 4
	MenuStartX       = 76
	MenuStartY       = 42
	MenuEntrySpacing = 50
)

var ErrMenuEmpty = errors.New("empty menu")

type Menu struct {
	stdImg         *ebiten.Image
	stdXSh, stdYSh float64

	current  *preset.Node
	selected int

	entries map[*preset.Node]*Entry
}

func NewMenu(asts *assets.Loader, menu *preset.Node) (*Menu, error) {
	if len(menu.Children) == 0 {
		return nil, ErrMenuEmpty
	}

	stdImg, err := asts.GetImage("ui/selected")
	if err != nil {
		return nil, err
	}

	m := &Menu{
		current: menu,
		stdImg:  stdImg,
	}

	err = m.buildEntries(asts, menu)
	if err != nil {
		return nil, err
	}

	sbds := m.stdImg.Bounds()
	ebds := m.entries[menu.Children[0]].Bounds() // arbitrary entry

	m.stdXSh = float64(sbds.Dx()-ebds.Dx()) / 2
	m.stdYSh = float64(sbds.Dy()-ebds.Dy()) / 2

	return m, nil
}

func (m *Menu) Draw(screen *ebiten.Image) {

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(MenuStartX, MenuStartY)

	for i := m.selected - 1; i < m.selected+MenuItems-1; i++ {
		p := (i + len(m.current.Children)) % len(m.current.Children)
		entry := m.entries[m.current.Children[p]]

		opts.GeoM.Translate(0, MenuEntrySpacing)
		screen.DrawImage(entry.Image, opts)

		if i == m.selected {
			// Draw stdImg background
			selOpts := &ebiten.DrawImageOptions{}
			selOpts.GeoM.Concat(opts.GeoM)
			selOpts.GeoM.Translate(-m.stdXSh, -m.stdYSh)
			screen.DrawImage(m.stdImg, selOpts)
		}
	}
}

func (m *Menu) Scroll(delta int) {
	m.selected -= delta
	if m.selected < 0 {
		m.selected = len(m.current.Children) - 1
	}
}

func (m *Menu) Forward() {

}

func (m *Menu) Backward() {

}

func (m *Menu) buildEntries(asts *assets.Loader, node *preset.Node) error {
	if m.entries == nil {
		m.entries = make(map[*preset.Node]*Entry)
	}

	for _, ch := range node.Children {
		et, err := NewEntry(asts, ch.Label)
		if err != nil {
			return err
		}
		m.entries[ch] = et
		if len(ch.Children) > 0 {
			if err = m.buildEntries(asts, ch); err != nil {
				return err
			}
		}
	}

	return nil
}
