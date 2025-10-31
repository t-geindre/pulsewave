package ui

import (
	"errors"
	"math"
	"synth/assets"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
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

	// Animation
	cursorY         float64
	targetCursorY   float64
	windowOffset    float64
	targetWinOffset float64
	animatingWin    bool
	animT           float64 // interpolation progress [0..1]
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
		current:   menu,
		cursorImg: cursorImg,
	}

	if err := m.buildEntries(asts, menu); err != nil {
		return nil, err
	}

	sbds := m.cursorImg.Bounds()
	ebds := m.entries[menu.Children[0]].Bounds()
	m.cursorXSh = float64(sbds.Dx()-ebds.Dx()) / 2
	m.cursorYSh = float64(sbds.Dy()-ebds.Dy()) / 2

	m.menuWindow = make([]int, MenuItems+2)
	for i := range m.menuWindow {
		m.menuWindow[i] = i % len(menu.Children)
	}

	m.cursorPos = 0
	m.cursorY = float64(m.cursorPos+1) * MenuEntrySpacing
	m.targetCursorY = m.cursorY

	return m, nil
}

// --- easing cubic: easeOutCubic ---
func easeOutCubic(t float64) float64 {
	return 1 - math.Pow(1-t, 3)
}

func (m *Menu) Update() {
	const speed = 0.22 // plus petit = plus lent

	if m.animatingWin {
		m.animT += speed
		if m.animT >= 1 {
			m.animT = 1
			m.finishScroll()
		}

		// easing entre 0 et targetWinOffset
		e := easeOutCubic(m.animT)
		m.windowOffset = m.targetWinOffset * e
	} else {
		// anime juste le curseur
		m.cursorY += (m.targetCursorY - m.cursorY) * 0.4
	}
}

func (m *Menu) finishScroll() {
	m.windowOffset = 0
	m.animatingWin = false
	m.animT = 0

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

func (m *Menu) Draw(screen *ebiten.Image) {
	for i, idx := range m.menuWindow {
		entry := m.entries[m.current.Children[idx]]
		y := MenuStartY + float64(i)*MenuEntrySpacing + m.windowOffset

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(MenuStartX, y)
		screen.DrawImage(entry.Image, opts)
	}

	// curseur à position interpolée
	selOpts := &ebiten.DrawImageOptions{}
	selY := MenuStartY + m.cursorY - m.cursorYSh
	selOpts.GeoM.Translate(MenuStartX-m.cursorXSh, selY)
	screen.DrawImage(m.cursorImg, selOpts)
}

func (m *Menu) Scroll(delta int) {
	if m.animatingWin || delta == 0 {
		return // ignore si animation en cours
	}

	total := len(m.current.Children)

	if delta < 0 {
		if m.cursorPos > 0 {
			m.cursorPos--
			m.targetCursorY = float64(m.cursorPos+1) * MenuEntrySpacing
		} else {
			// scroll la fenêtre vers le haut (le curseur monte, donc le contenu descend)
			m.animatingWin = true
			m.targetWinOffset = +MenuEntrySpacing
			m.scrollingUp = true
		}
	}

	if delta > 0 {
		if m.cursorPos < MenuItems-1 {
			m.cursorPos++
			m.targetCursorY = float64(m.cursorPos+1) * MenuEntrySpacing
		} else {
			// scroll la fenêtre vers le bas (le curseur descend, donc le contenu monte)
			m.animatingWin = true
			m.targetWinOffset = -MenuEntrySpacing
			m.scrollingDown = true
		}
	}

	m.cursorPos = (m.cursorPos + total) % total
}

// --- Garde la construction récursive des entrées ---
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
