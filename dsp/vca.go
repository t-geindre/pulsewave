package dsp

import "synth/audio"

type Vca struct {
	Node
	gain Param
}

func NewVca(src Node, gain Param) *Vca {
	return &Vca{
		Node: src,
		gain: gain,
	}
}

func (v *Vca) Process(b *audio.Block) {
	v.Node.Process(b)

	g := v.gain.Resolve(b.Cycle)
	for i := 0; i < audio.BlockSize; i++ {
		b.L[i] *= g[i]
		b.R[i] *= g[i]
	}
}

func (v *Vca) Reset() {
	v.Node.Reset()
}
