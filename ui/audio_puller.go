package ui

import "synth/dsp"

type AudioPuller struct {
	dsp.Node
	out *AudioQueue
}

func NewAudioPuller(src dsp.Node, out *AudioQueue) *AudioPuller {
	return &AudioPuller{
		Node: src,
		out:  out,
	}
}

func (a *AudioPuller) Process(block *dsp.Block) {
	a.Node.Process(block)
	a.out.TryWrite(*block)
}
