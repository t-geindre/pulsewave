package ui

import (
	"synth/assets"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
)

type Selector struct {
	bg   *ebiten.Image
	node *preset.SelectorNode
}

func NewSelector(asts *assets.Loader, node *preset.SelectorNode) (*Selector, error) {
	bg, err := asts.GetImage("ui/selector/bg")
	if err != nil {
		return nil, err
	}

	return &Selector{
		bg: bg,
	}, nil
}

func (s *Selector) Draw(image *ebiten.Image) {
	image.DrawImage(s.bg, nil)
}

func (s *Selector) Update() {
}

func (s *Selector) Scroll(delta int) {
}

func (s *Selector) CurrentTarget() preset.Node {
	return nil
}
