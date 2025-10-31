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
	cursorImg            *ebiten.Image
	cursorXSh, cursorYSh float64

	menuWindow []int
	cursorPos  int

	current *preset.Node

	entries map[*preset.Node]*Entry
}

func NewMenu(asts *assets.Loader, menu *preset.Node) (*Menu, error) {
	if len(menu.Children) == 0 {
		return nil, ErrMenuEmpty
	}

	cursorImg, err := asts.GetImage("ui/selected")
	if err != nil {
		return nil, err
	}

	m := &Menu{
		current:   menu,
		cursorImg: cursorImg,
	}

	err = m.buildEntries(asts, menu)
	if err != nil {
		return nil, err
	}

	sbds := m.cursorImg.Bounds()
	ebds := m.entries[menu.Children[0]].Bounds() // arbitrary entry

	m.cursorXSh = float64(sbds.Dx()-ebds.Dx()) / 2
	m.cursorYSh = float64(sbds.Dy()-ebds.Dy()) / 2

	m.menuWindow = make([]int, MenuItems+2)
	for i := range m.menuWindow {
		m.menuWindow[i] = i % len(menu.Children)
	}

	return m, nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	for i, idx := range m.menuWindow {
		entry := m.entries[m.current.Children[idx]]

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(MenuStartX, MenuStartY+float64(i)*MenuEntrySpacing)
		screen.DrawImage(entry.Image, opts)
		if i == m.cursorPos+1 {
			selOpts := &ebiten.DrawImageOptions{}
			selOpts.GeoM.Concat(opts.GeoM)
			selOpts.GeoM.Translate(-m.cursorXSh, -m.cursorYSh)
			screen.DrawImage(m.cursorImg, selOpts)
		}
	}
}

func (m *Menu) Scroll(delta int) {
	total := len(m.current.Children)
	window := len(m.menuWindow)

	if delta < 0 {
		if m.cursorPos > 0 {
			m.cursorPos--
		} else {
			for i := window - 1; i > 0; i-- {
				m.menuWindow[i] = m.menuWindow[i-1]
			}
			m.menuWindow[0] = (m.menuWindow[0] - 1 + total) % total
		}
	}

	if delta > 0 {
		if m.cursorPos < MenuItems-1 {
			m.cursorPos++
		} else {
			for i := 0; i < window-1; i++ {
				m.menuWindow[i] = m.menuWindow[i+1]
			}
			m.menuWindow[window-1] = (m.menuWindow[window-1] + 1) % total
		}
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
