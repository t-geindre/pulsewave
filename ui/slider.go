package ui

import (
	"synth/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Todo get is from config
const (
	SliderBoxStartX          = 81 - 51
	SliderBoxStartY          = 107 - 70
	SliderBoxWidth           = 349 - 57
	SliderBoxHeight          = 236 - 11
	SliderValueBottomSpacing = 140
	SliderTitleTopSpacing    = 10
	SliderBackLabelX         = 95 - 51
	SliderBackLabelY         = 245 - 70
)

type Slider struct {
	bg        *ebiten.Image
	faceParam text.Face
	faceBack  text.Face
	node      *ParameterNode
}

func NewSlider(asts *assets.Loader, sn *ParameterNode) (*Slider, error) {
	bg, err := asts.GetImage("ui/slider/bg")
	if err != nil {
		return nil, err
	}

	faceParam, err := asts.GetFace("ui/param")
	if err != nil {
		return nil, err
	}

	faceBack, err := asts.GetFace("ui/param_back")
	if err != nil {
		return nil, err
	}

	return &Slider{
		bg:        bg,
		faceParam: faceParam,
		faceBack:  faceBack,
		node:      sn,
	}, nil
}

func (s *Slider) Draw(image *ebiten.Image) {
	image.DrawImage(s.bg, nil)

	// Draw Value
	valueOpts := &text.DrawOptions{}
	valueOpts.GeoM.Translate(SliderBoxStartX, SliderBoxStartY+SliderBoxHeight)

	vDisplay := s.node.Display()
	ww, wh := text.Measure(vDisplay, s.faceParam, 0)
	valueOpts.GeoM.Translate((SliderBoxWidth-ww)/2, -wh/2-SliderValueBottomSpacing)

	text.Draw(image, vDisplay, s.faceParam, valueOpts)

	// Draw back label
	if p := s.node.Parent(); p != nil {
		opt := &text.DrawOptions{}
		opt.GeoM.Translate(SliderBackLabelX, SliderBackLabelY)
		text.Draw(image, p.Label(), s.faceBack, opt)
	}
}

func (s *Slider) Update() {
}

func (s *Slider) Scroll(delta int) {
	s.node.SetVal(s.node.Val() + s.node.Step()*float32(-delta))
}

func (s *Slider) CurrentTarget() Node {
	return nil
}
