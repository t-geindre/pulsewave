package preset

import "synth/dsp"

type Skipper struct {
	dsp.Param
	normal, skipped dsp.Node
}

func NewSkipper(normal, skipped dsp.Node, p dsp.Param) *Skipper {
	return &Skipper{
		Param:   p,
		normal:  normal,
		skipped: skipped,
	}
}

func (s *Skipper) Process(block *dsp.Block) {
	if s.Resolve(block.Cycle)[0] < 0.5 {
		s.skipped.Process(block)
		return
	}

	s.normal.Process(block)
}

func (s *Skipper) Reset() {
	s.normal.Reset()
}
