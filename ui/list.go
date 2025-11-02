package ui

import (
	"errors"
	"synth/assets"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
)

// Todo get it from config
const (
	ListVisibleItems = 4
	ListPaddingTop   = 20
	ListPaddingLeft  = 25
	ListEntrySpacing = 50
)

var ErrEmptyList = errors.New("empty component")

type List struct {
	node    *preset.Node
	entries map[*preset.Node]*ListEntry

	listWindow []int
	cursorPos  int

	cursorImg            *ebiten.Image
	cursorXSh, cursorYSh float64

	loop bool

	cursorY         float64
	targetCursorY   float64
	windowOffset    float64
	targetWinOffset float64
	animatingWin    bool
	animT           float64
	scrollingUp     bool
	scrollingDown   bool
}

func NewList(asts *assets.Loader, node *preset.Node) (*List, error) {
	if len(node.Children) == 0 {
		return nil, ErrEmptyList
	}

	cursorImg, err := asts.GetImage("ui/selected")
	if err != nil {
		return nil, err
	}

	l := &List{
		node:      node,
		cursorImg: cursorImg,
		entries:   make(map[*preset.Node]*ListEntry),
	}

	if err = l.buildEntries(asts, node); err != nil {
		return nil, err
	}

	// Cursor alignment
	sbds := l.cursorImg.Bounds()
	ebds := l.entries[node.Children[0]].Bounds() // arbitrary entry

	l.cursorXSh = -float64(sbds.Dx()-ebds.Dx()) / 2
	l.cursorYSh = float64(sbds.Dy()-ebds.Dy()) / 2

	// List window + Loop mode
	l.loop = true
	ws := ListVisibleItems + 2
	if len(node.Children) <= ListVisibleItems {
		ws = len(node.Children)
		l.loop = false
	}
	l.listWindow = make([]int, ws)
	for i := range l.listWindow {
		s := i
		if l.loop {
			s--
			if s < 0 {
				s = len(l.listWindow) - 1
			}
		}
		l.listWindow[i] = s % len(node.Children)
	}

	// Init pos
	l.cursorPos = 0
	l.cursorY = float64(l.cursorPos)*ListEntrySpacing + ListPaddingTop
	l.targetCursorY = l.cursorY

	return l, nil
}

func (l *List) Update() {
	const speed = 0.22 // todo move to config
	if l.animatingWin {
		l.animT += speed
		if l.animT >= 1 {
			l.animT = 1
			l.finishScroll()
		}
		l.windowOffset = l.targetWinOffset * easeOutCubic(l.animT)
	} else {
		l.cursorY += (l.targetCursorY - l.cursorY) * 0.4 // todo move const val to config
	}
}

func (l *List) Draw(screen *ebiten.Image) {
	// Entries
	for i, idx := range l.listWindow {
		entry := l.entries[l.node.Children[idx]]
		startIndex := 0
		if l.loop {
			startIndex = -1
		}
		y := float64(i+startIndex)*ListEntrySpacing + l.windowOffset + ListPaddingTop

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(ListPaddingLeft, y)
		screen.DrawImage(entry.Image, opts)
	}

	// Cursor
	selOpts := &ebiten.DrawImageOptions{}
	selY := l.cursorY - l.cursorYSh
	selOpts.GeoM.Translate(l.cursorXSh+ListPaddingLeft, selY)
	screen.DrawImage(l.cursorImg, selOpts)
}

func (l *List) Scroll(delta int) {
	if l.animatingWin || delta == 0 {
		return
	}

	if !l.loop {
		switch {
		case delta < 0:
			if l.cursorPos > 0 {
				l.moveCursor(-1)
			}
		default:
			if l.cursorPos < len(l.listWindow)-1 {
				l.moveCursor(1)
			}
		}
		return
	}

	total := len(l.node.Children)

	switch {
	case delta < 0:
		if l.cursorPos > 0 {
			l.moveCursor(-1)
		} else {
			l.startScroll(+ListEntrySpacing, true)
		}
	default:
		if l.cursorPos < ListVisibleItems-1 {
			l.moveCursor(+1)
		} else {
			l.startScroll(-ListEntrySpacing, false)
		}
	}

	l.cursorPos = (l.cursorPos + total) % total
}

func (l *List) CurrentTarget() *preset.Node {
	pos := l.cursorPos
	if l.loop {
		pos++
	}
	return l.node.Children[l.listWindow[pos]]
}

func (l *List) moveCursor(dir int) {
	l.cursorPos += dir
	l.targetCursorY = float64(l.cursorPos)*ListEntrySpacing + ListPaddingTop
}

func (l *List) startScroll(offset float64, up bool) {
	l.animatingWin = true
	l.targetWinOffset = offset
	l.scrollingUp = up
	l.scrollingDown = !up
	l.animT = 0
}

func (l *List) finishScroll() {
	l.windowOffset, l.animT, l.animatingWin = 0, 0, false

	total := len(l.node.Children)
	window := len(l.listWindow)

	if l.scrollingDown {
		for i := 0; i < window-1; i++ {
			l.listWindow[i] = l.listWindow[i+1]
		}
		l.listWindow[window-1] = (l.listWindow[window-1] + 1) % total
		l.scrollingDown = false
	}
	if l.scrollingUp {
		for i := window - 1; i > 0; i-- {
			l.listWindow[i] = l.listWindow[i-1]
		}
		l.listWindow[0] = (l.listWindow[0] - 1 + total) % total
		l.scrollingUp = false
	}
}
func (l *List) buildEntries(asts *assets.Loader, node *preset.Node) error {
	for _, ch := range node.Children {
		entry, err := NewListEntry(asts, ch.Label)
		if err != nil {
			return err
		}
		l.entries[ch] = entry

	}
	return nil
}
