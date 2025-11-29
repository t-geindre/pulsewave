package ui

import (
	"synth/assets"
	"synth/tree"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// todo get it from config
const (
	SelectorBackLabelX = 44
	SelectorBackLabelY = 175
	BoxStartX          = 38
	BoxStartY          = 40
	BoxWith            = 297
	BoxHeight          = 126
	IconSpacing        = 25
	ScrollSpeed        = 10
)

type Selector struct {
	bg          *ebiten.Image
	node        tree.SelectorNode
	faceBack    text.Face
	faceOption  text.Face
	icons       map[float32]*ebiten.Image
	forward     *ebiten.Image
	clippingBox *ebiten.Image

	currentIndex int
	prevIndex    int

	scrollY       float64
	targetY       float64
	isMoving      bool
	dir           int
	currentOffset float64
}

func NewSelector(asts *assets.Loader, node tree.SelectorNode) (*Selector, error) {
	bg, err := asts.GetImage("ui/selector/bg")
	if err != nil {
		return nil, err
	}

	faceBack, err := asts.GetFace("ui/selector/back")
	if err != nil {
		return nil, err
	}

	faceOption, err := asts.GetFace("ui/selector/option")
	if err != nil {
		return nil, err
	}

	var forward *ebiten.Image
	if node.RequiresValidation() {
		forward, err = asts.GetImage("ui/arrow_froward")
		if err != nil {
			return nil, err
		}
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
		faceOption:  faceOption,
		forward:     forward,
		icons:       icons,
		clippingBox: ebiten.NewImage(BoxWith, BoxHeight),
		prevIndex:   0,
	}, nil
}

func (s *Selector) Update() {
	if !s.isMoving {
		return
	}

	const duration = BoxHeight / ScrollSpeed
	if s.scrollY < 1.0 {
		s.scrollY += 1.0 / duration
		if s.scrollY > 1.0 {
			s.scrollY = 1.0
		}
	}

	t := easeOutCubic(s.scrollY)
	currentOffset := t * s.targetY

	if s.scrollY >= 1.0 {
		s.isMoving = false
		s.scrollY = 0
		s.targetY = 0
	}

	s.currentOffset = currentOffset
}

func (s *Selector) Scroll(delta int) {
	if s.isMoving {
		return
	}

	delta = -delta
	if delta > 1 {
		delta = 1
	} else if delta < -1 {
		delta = -1
	}

	prev := s.currentIndex
	s.currentIndex = prev + delta
	for s.currentIndex < 0 {
		s.currentIndex += len(s.node.Options())
	}
	s.currentIndex = s.currentIndex % len(s.node.Options())

	s.node.SetVal(s.node.Options()[s.currentIndex].Value())
	s.prevIndex = prev
	s.isMoving = true

	if delta > 0 {
		s.dir = 1
		s.targetY = BoxHeight
	} else {
		s.dir = -1
		s.targetY = -BoxHeight
	}
}

func (s *Selector) Draw(image *ebiten.Image) {
	image.DrawImage(s.bg, nil)
	s.clippingBox.Clear()

	options := s.node.Options()

	if !s.isMoving {
		s.drawOption(s.clippingBox, options[s.currentIndex], 0)
	} else {
		prev := s.prevIndex
		dir := s.dir
		offset := s.currentOffset
		s.drawOption(s.clippingBox, options[prev], -offset)
		s.drawOption(s.clippingBox, options[s.currentIndex], float64(dir)*BoxHeight-offset)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(BoxStartX, BoxStartY)
	image.DrawImage(s.clippingBox, op)

	if p := s.node.Parent(); p != nil {
		opt := &text.DrawOptions{}
		opt.GeoM.Translate(SelectorBackLabelX, SelectorBackLabelY)
		text.Draw(image, p.Label(), s.faceBack, opt)
	}
}

func (s *Selector) drawOption(dst *ebiten.Image, opt *tree.SelectorOption, y float64) {
	icon, hasIcon := s.icons[opt.Value()]
	tw, th := text.Measure(opt.Label(), s.faceOption, 0)

	var iconW, iconH float64
	if hasIcon && icon != nil {
		bds := icon.Bounds()
		iconW = float64(bds.Dx())
		iconH = float64(bds.Dy())
	}

	spacing := 0.0
	if hasIcon && icon != nil {
		spacing = IconSpacing
	}

	width := iconW + spacing + tw
	startX := (BoxWith - width) / 2

	centerY := float64(BoxHeight) / 2

	if hasIcon && icon != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(startX, centerY-iconH/2+y)
		dst.DrawImage(icon, op)
	}

	txt := &text.DrawOptions{}
	txtX := startX + iconW + spacing
	txtY := centerY - th/2 + y
	txt.GeoM.Translate(txtX, txtY)
	text.Draw(dst, opt.Label(), s.faceOption, txt)

	if s.forward != nil {
		op := &ebiten.DrawImageOptions{}
		bds := s.forward.Bounds()
		fw := float64(bds.Dx())
		fh := float64(bds.Dy())
		op.GeoM.Translate(BoxWith/2-fw/2+width, centerY-fh/2+y)
		dst.DrawImage(s.forward, op)
	}
}

func (s *Selector) CurrentTarget() tree.Node {
	if s.node.RequiresValidation() {
		s.node.Validate()
		return s.node.Parent()
	}

	return nil
}

func (s *Selector) Focus() {
	s.currentIndex = 0
	v := s.node.Val()
	for i, opt := range s.node.Options() {
		if opt.Value() == v {
			s.currentIndex = i
			break
		}
	}
}
func (s *Selector) Blur() {}
