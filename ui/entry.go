package ui

import (
	"synth/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Todo get it from config
const (
	LetterSpacing = 2
	EntryWidth    = 329.0
	EntryHeight   = 24.0
)

type Entry struct {
	*ebiten.Image
}

func NewEntry(asts *assets.Loader, str string) (*Entry, error) {
	w, h := EntryWidth, EntryHeight
	img := ebiten.NewImage(int(w), int(h))

	face, err := asts.GetFace("ui/face")
	if err != nil {
		return nil, err
	}

	opts := &text.DrawOptions{}
	_, th := text.Measure(str, face, 0)
	opts.GeoM.Translate(0, (h-th)/2)

	for _, r := range str {
		s := string(r)
		text.Draw(img, s, face, opts)
		opts.GeoM.Translate(text.Advance(s, face)+LetterSpacing, 0)
	}

	arrow, err := asts.GetImage("ui/arrow")
	if err != nil {
		return nil, err
	}

	bds := arrow.Bounds()

	arrOpts := &ebiten.DrawImageOptions{}
	arrOpts.GeoM.Translate(w-float64(bds.Dx()), (h-float64(bds.Dy()))/2)

	img.DrawImage(arrow, arrOpts)

	return &Entry{
		Image: img,
	}, nil
}
