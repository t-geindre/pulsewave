package ui

import (
	"synth/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Slider struct {
	bg   *ebiten.Image
	face text.Face
	node *ParameterNode
}

func NewSlider(asts *assets.Loader, sn *ParameterNode) (*Slider, error) {
	bg, err := asts.GetImage("ui/slider/bg")
	if err != nil {
		return nil, err
	}

	face, err := asts.GetFace("ui/face")
	if err != nil {
		return nil, err
	}

	return &Slider{
		bg:   bg,
		face: face,
		node: sn,
	}, nil
}

func (s *Slider) Draw(image *ebiten.Image) {
	image.DrawImage(s.bg, nil)
	text.Draw(image, s.node.Display(), s.face, nil)
}

func (s *Slider) Update() {
}

func (s *Slider) Scroll(delta int) {
	s.node.SetVal(s.node.Val() + s.node.Step()*float32(-delta))
}

func (s *Slider) CurrentTarget() Node {
	return nil
}
