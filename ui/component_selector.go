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
	BoxStartX          = 38
	BoxStartY          = 46
	BoxWith            = 297
	BoxHeight          = 114
	IconSpacing        = 25
)

type Selector struct {
	bg          *ebiten.Image
	node        *preset.SelectorNode
	faceBack    text.Face
	icons       map[float32]*ebiten.Image
	clippingBox *ebiten.Image
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

	icons := make(map[float32]*ebiten.Image)
	for _, opt := range node.Options() {
		if opt.Icon() != "" {
			iconImg, err := asts.GetImage(opt.Icon())
			if err != nil {
				return nil, err
			}
			icons[opt.Value()] = iconImg
		}
	}

	return &Selector{
		bg:          bg,
		node:        node,
		faceBack:    faceBack,
		icons:       icons,
		clippingBox: ebiten.NewImage(BoxWith, BoxHeight),
	}, nil
}

func (s *Selector) Draw(image *ebiten.Image) {
	image.DrawImage(s.bg, nil)
	s.clippingBox.Clear()

	// Options
	cur := int(s.node.Val())

	for _, idx := range []int{
		cur,
	} {
		option := s.node.Options()[idx]
		icon := s.icons[option.Value()]

		tw, th := text.Measure(option.Label(), s.faceBack, 0)
		iconBds := icon.Bounds()

		width := float64(iconBds.Dx()) + tw + IconSpacing
		startX := (BoxWith - width) / 2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(startX, (float64(BoxHeight)-float64(iconBds.Dy()))/2)
		s.clippingBox.DrawImage(icon, op)

		opt := &text.DrawOptions{}
		opt.GeoM.Translate(startX+float64(iconBds.Dx())+IconSpacing, float64(iconBds.Dy())/2+float64(th)/2)
		text.Draw(s.clippingBox, option.Label(), s.faceBack, opt)
	}

	// Clipping box draw
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(BoxStartX, BoxStartY)
	image.DrawImage(s.clippingBox, op)

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
