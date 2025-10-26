package ui

import (
	"synth/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Entry struct {
	*ebiten.Image
}

func NewEntry(asts *assets.Loader, str string) *Entry {
	const Spacing = 2 // Todo get it from config
	const Width = 329.0
	const Height = 24.0

	w, h := Width, Height
	img := ebiten.NewImage(int(w), int(h))

	face := asts.GetFace("ui/face")

	opts := &text.DrawOptions{}
	_, th := text.Measure(str, face, 0)
	opts.GeoM.Translate(0, (h-th)/2)

	for _, r := range str {
		s := string(r)
		text.Draw(img, s, face, opts)
		opts.GeoM.Translate(text.Advance(s, face)+Spacing, 0)
	}

	arrow := asts.GetImage("ui/arrow")
	bds := arrow.Bounds()

	arrOpts := &ebiten.DrawImageOptions{}
	arrOpts.GeoM.Translate(w-float64(bds.Dx()), (h-float64(bds.Dy()))/2)

	img.DrawImage(arrow, arrOpts)

	return &Entry{
		Image: img,
	}
}
