package ui

import (
	"synth/assets"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// todo get it from config
const (
	SelectorBackLabelX = 44
	SelectorBackLabelY = 175
)

type Selector struct {
	bg       *ebiten.Image
	node     *preset.SelectorNode
	faceBack text.Face
}

func NewSelector(asts *assets.Loader, node *preset.SelectorNode) (*Selector, error) {
	bg, err := asts.GetImage("ui/selector/bg")
	if err != nil {
		return nil, err
	}

	faceBack, err := asts.GetFace("ui/selector/back")
	if err != nil {
		return nil, err
	}

	return &Selector{
		bg:       bg,
		node:     node,
		faceBack: faceBack,
	}, nil
}

func (s *Selector) Draw(image *ebiten.Image) {
	image.DrawImage(s.bg, nil)

	// Current
	idx := int(s.node.Val())
	text.Draw(image, s.node.Options()[idx].Label(), s.faceBack, nil)

	// Draw back label
	if p := s.node.Parent(); p != nil {
		opt := &text.DrawOptions{}
		opt.GeoM.Translate(SelectorBackLabelX, SelectorBackLabelY)
		text.Draw(image, p.Label(), s.faceBack, opt)
	}
}

func (s *Selector) Update() {
}

func (s *Selector) Scroll(delta int) {
	v := int(s.node.Val()) + delta
	for v < 0 {
		v += len(s.node.Options())
	}
	v = v % len(s.node.Options())
	s.node.SetVal(float32(v))
}

func (s *Selector) CurrentTarget() preset.Node {
	return nil
}
