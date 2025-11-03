package ui

import (
	"synth/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Slider struct {
	bg *ebiten.Image
}

func NewSlider(asts *assets.Loader, sn *SliderNode) (*Slider, error) {
	bg, err := asts.GetImage("ui/slider/bg")
	if err != nil {
		return nil, err
	}

	return &Slider{
		bg: bg,
	}, nil
}

func (s *Slider) Draw(image *ebiten.Image) {
	image.DrawImage(s.bg, nil)
}

func (s *Slider) Update() {
}

func (s *Slider) Scroll(delta int) {
}

func (s *Slider) CurrentTarget() Node {
	return nil
}
