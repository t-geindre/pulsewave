package ui

import (
	"errors"
	"synth/assets"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
)

// Todo get it from config
const (
	MenuItems        = 4
	MenuHeight       = 211
	MenuWidth        = 376
	MenuStartX       = 51
	MenuStartY       = 70
	MenuPaddingTop   = 20
	MenuPaddingLeft  = 25
	MenuEntrySpacing = 50
)

var ErrMenuEmpty = errors.New("empty menu")

type Menu struct {
	current *preset.Node
	entries map[*preset.Node]*Entry

	menuWindow []int
	cursorPos  int

	cursorImg            *ebiten.Image
	cursorXSh, cursorYSh float64
	clippingMask         *ebiten.Image

	cursorY         float64
	targetCursorY   float64
	windowOffset    float64
	targetWinOffset float64
	animatingWin    bool
	animT           float64
	scrollingUp     bool
	scrollingDown   bool
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
		current:      menu,
		cursorImg:    cursorImg,
		entries:      make(map[*preset.Node]*Entry),
		clippingMask: ebiten.NewImage(MenuWidth, MenuHeight),
	}

	if err = m.buildEntries(asts, menu); err != nil {
		return nil, err
	}

	// Cursor alignment
	sbds := m.cursorImg.Bounds()
	ebds := m.entries[menu.Children[0]].Bounds() // arbitrary entry
	m.cursorXSh = float64(sbds.Dx()-ebds.Dx()) / 2
	m.cursorYSh = float64(sbds.Dy()-ebds.Dy()) / 2

	// Menu window
	m.menuWindow = make([]int, MenuItems+2)
	for i := range m.menuWindow {
		m.menuWindow[i] = i % len(menu.Children)
	}

	// Init pos
	m.cursorPos = 0
	m.cursorY = float64(m.cursorPos)*MenuEntrySpacing + MenuPaddingTop
	m.targetCursorY = m.cursorY

	return m, nil
}

func (m *Menu) Update() {
	const speed = 0.22
	if m.animatingWin {
		m.animT += speed
		if m.animT >= 1 {
			m.animT = 1
			m.finishScroll()
		}
		m.windowOffset = m.targetWinOffset * easeOutCubic(m.animT)
	} else {
		m.cursorY += (m.targetCursorY - m.cursorY) * 0.4
	}
}

func (m *Menu) Draw(screen *ebiten.Image) {
	m.clippingMask.Clear()

	// Entries
	for i, idx := range m.menuWindow {
		entry := m.entries[m.current.Children[idx]]
		y := float64(i-1)*MenuEntrySpacing + m.windowOffset + MenuPaddingTop

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(MenuPaddingLeft, y)
		m.clippingMask.DrawImage(entry.Image, opts)
	}

	// Clip mask rendering
	menuOpts := &ebiten.DrawImageOptions{}
	menuOpts.GeoM.Translate(MenuStartX, MenuStartY)
	screen.DrawImage(m.clippingMask, menuOpts)

	// Cursor
	selOpts := &ebiten.DrawImageOptions{}
	selY := MenuStartY + m.cursorY - m.cursorYSh
	selOpts.GeoM.Translate(MenuStartX-m.cursorXSh+MenuPaddingLeft, selY)
	screen.DrawImage(m.cursorImg, selOpts)
}

func (m *Menu) Scroll(delta int) {
	if m.animatingWin || delta == 0 {
		return
	}

	total := len(m.current.Children)

	switch {
	case delta < 0:
		if m.cursorPos > 0 {
			m.moveCursor(-1)
		} else {
			m.startScroll(+MenuEntrySpacing, true)
		}
	default:
		if m.cursorPos < MenuItems-1 {
			m.moveCursor(+1)
		} else {
			m.startScroll(-MenuEntrySpacing, false)
		}
	}

	m.cursorPos = (m.cursorPos + total) % total
}

func (m *Menu) moveCursor(dir int) {
	m.cursorPos += dir
	m.targetCursorY = float64(m.cursorPos)*MenuEntrySpacing + MenuPaddingTop
}

func (m *Menu) startScroll(offset float64, up bool) {
	m.animatingWin = true
	m.targetWinOffset = offset
	m.scrollingUp = up
	m.scrollingDown = !up
	m.animT = 0
}

func (m *Menu) finishScroll() {
	m.windowOffset, m.animT, m.animatingWin = 0, 0, false

	total := len(m.current.Children)
	window := len(m.menuWindow)

	if m.scrollingDown {
		for i := 0; i < window-1; i++ {
			m.menuWindow[i] = m.menuWindow[i+1]
		}
		m.menuWindow[window-1] = (m.menuWindow[window-1] + 1) % total
		m.scrollingDown = false
	}
	if m.scrollingUp {
		for i := window - 1; i > 0; i-- {
			m.menuWindow[i] = m.menuWindow[i-1]
		}
		m.menuWindow[0] = (m.menuWindow[0] - 1 + total) % total
		m.scrollingUp = false
	}
}
func (m *Menu) buildEntries(asts *assets.Loader, node *preset.Node) error {
	for _, ch := range node.Children {
		entry, err := NewEntry(asts, ch.Label)
		if err != nil {
			return err
		}
		m.entries[ch] = entry

	}
	return nil
}
