package ui

import (
	"errors"
	"synth/assets"
	"synth/tree"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Todo get it from config
const (
	ListVisibleItems   = 4
	ListPaddingTop     = 20
	ListHorPadding     = 25
	ListEntrySpacing   = 50
	ListEntryWidth     = 329.0
	ListEntryHeight    = 24.0
	ListPreviewSpacing = 10
)

var ErrEmptyList = errors.New("empty list")

type List struct {
	node tree.Node

	firstIndex int
	visible    int
	cursorPos  int

	cursorImg            *ebiten.Image
	cursorXSh, cursorYSh float64

	faceEntry  text.Face
	arrowEntry *ebiten.Image

	cursorY         float64
	targetCursorY   float64
	windowOffset    float64
	targetWinOffset float64
	animatingWin    bool
	animT           float64
	scrollingUp     bool
	scrollingDown   bool
}

func NewList(asts *assets.Loader, node tree.Node) (*List, error) {
	if len(node.Children()) == 0 {
		return nil, ErrEmptyList
	}

	cursorImg, err := asts.GetImage("ui/list/selected")
	if err != nil {
		return nil, err
	}

	faceEntry, err := asts.GetFace("ui/list/entry")
	if err != nil {
		return nil, err
	}

	arrowEntry, err := asts.GetImage("ui/arrow_froward")
	if err != nil {
		return nil, err
	}

	l := &List{
		node:       node,
		cursorImg:  cursorImg,
		faceEntry:  faceEntry,
		arrowEntry: arrowEntry,
	}

	// Cursor alignment
	children := node.Children()
	sbds := l.cursorImg.Bounds()

	l.cursorXSh = -float64(sbds.Dx()-ListEntryWidth) / 2
	l.cursorYSh = float64(sbds.Dy()-ListEntryHeight) / 2

	total := len(children)
	visible := ListVisibleItems
	if total < visible {
		visible = total
	}
	l.visible = visible
	l.firstIndex = 0

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
	total := len(l.node.Children())

	for i := 0; i < l.visible; i++ {
		idx := l.firstIndex + i
		if idx >= total {
			break
		}
		y := float64(i)*ListEntrySpacing + l.windowOffset + ListPaddingTop
		l.drawEntry(screen, idx, y)
	}

	if l.animatingWin {
		first := l.firstIndex

		if l.scrollingDown {
			bottomIdx := first + l.visible - 1
			if bottomIdx < total-1 {
				nextIdx := bottomIdx + 1
				y := float64(l.visible)*ListEntrySpacing + l.windowOffset + ListPaddingTop
				l.drawEntry(screen, nextIdx, y)
			}
		}

		if l.scrollingUp {
			if first > 0 {
				prevIdx := first - 1
				y := -1*ListEntrySpacing + l.windowOffset + ListPaddingTop
				l.drawEntry(screen, prevIdx, y)
			}
		}
	}

	// Cursor
	selOpts := &ebiten.DrawImageOptions{}
	selY := l.cursorY - l.cursorYSh
	selOpts.GeoM.Translate(l.cursorXSh+ListHorPadding, selY)
	screen.DrawImage(l.cursorImg, selOpts)
}

func (l *List) Scroll(delta int) {
	if l.animatingWin || delta == 0 {
		return
	}

	children := l.node.Children()
	total := len(children)
	if total == 0 || l.visible == 0 {
		return
	}

	first := l.firstIndex
	window := l.visible
	globalCursor := first + l.cursorPos
	lastGlobal := total - 1

	if delta < 0 {
		if globalCursor == 0 {
			return
		}

		if l.cursorPos > 0 {
			l.moveCursor(-1)
			return
		}

		if first > 0 {
			l.startScroll(+ListEntrySpacing, true)
		}
		return
	}

	if globalCursor == lastGlobal {
		return
	}

	if l.cursorPos < window-1 && first+l.cursorPos+1 < total {
		l.moveCursor(+1)
		return
	}

	bottomIndex := first + window - 1
	if bottomIndex < lastGlobal {
		l.startScroll(-ListEntrySpacing, false)
	}
}

func (l *List) CurrentTarget() tree.Node {
	children := l.node.Children()
	idx := l.firstIndex + l.cursorPos
	if idx < 0 || idx >= len(children) {
		return nil
	}
	return children[idx]
}

func (l *List) moveCursor(dir int) {
	l.cursorPos += dir
	if l.cursorPos < 0 {
		l.cursorPos = 0
	}
	if l.cursorPos >= l.visible {
		l.cursorPos = l.visible - 1
	}
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

	children := l.node.Children()
	total := len(children)
	window := l.visible
	if window == 0 || total == 0 {
		l.scrollingUp = false
		l.scrollingDown = false
		return
	}

	first := l.firstIndex

	if l.scrollingDown {
		lastAllowedFirst := total - window
		if first < lastAllowedFirst {
			first++
			l.firstIndex = first
		}
		l.scrollingDown = false
	}

	if l.scrollingUp {
		if first > 0 {
			first--
			l.firstIndex = first
		}
		l.scrollingUp = false
	}
}

func (l *List) drawEntry(screen *ebiten.Image, idx int, y float64) {
	entry := l.node.Children()[idx]

	// Label
	label := entry.Label()
	_, th := text.Measure(label, l.faceEntry, 0)
	textCenterY := (ListEntryHeight - th) / 2

	entry.Label()
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(ListHorPadding, y+textCenterY)
	text.Draw(screen, entry.Label(), l.faceEntry, opts)

	// Arrow
	arrBds := l.arrowEntry.Bounds()
	arrOpts := &ebiten.DrawImageOptions{}
	arrOpts.GeoM.Translate(ListEntryWidth-float64(arrBds.Dx())+ListHorPadding, y+(ListEntryHeight-float64(arrBds.Dy()))/2)
	screen.DrawImage(l.arrowEntry, arrOpts)

	// Preview
	preview := entry.Preview()
	if preview != "" {
		prewOpts := &text.DrawOptions{}
		pw, _ := text.Measure(preview, l.faceEntry, 0)
		prewOpts.GeoM.Translate(ListEntryWidth-pw-float64(arrBds.Dx())+ListHorPadding-ListPreviewSpacing, y+textCenterY)
		text.Draw(screen, preview, l.faceEntry, prewOpts)
	}
}
