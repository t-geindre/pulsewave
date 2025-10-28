package dsp

import "synth/audio"

type Vca struct {
	audio.Source
	gain Param
}

func NewVca(src audio.Source, gain Param) *Vca {
	return &Vca{
		Source: src,
		gain:   gain,
	}
}

func (v *Vca) Process(b *audio.Block) {
	v.Source.Process(b)

	g := v.gain.Resolve(b.Cycle)
	for i := 0; i < audio.BlockSize; i++ {
		b.L[i] *= g[i]
		b.R[i] *= g[i]
	}
}
