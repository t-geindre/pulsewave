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
}

func (s *Selector) CurrentTarget() preset.Node {
	return nil
}
