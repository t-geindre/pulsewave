package preset

import "synth/dsp"

type NodeSkipper struct {
	dsp.Param
	normal, skipped dsp.Node
}

func NewNodeSkipper(normal, skipped dsp.Node, toggle dsp.Param) *NodeSkipper {
	return &NodeSkipper{
		Param:   toggle,
		normal:  normal,
		skipped: skipped,
	}
}

func (s *NodeSkipper) Process(block *dsp.Block) {
	if s.Resolve(block.Cycle)[0] < 0.5 {
		s.skipped.Process(block)
		return
	}

	s.normal.Process(block)
}

func (s *NodeSkipper) Reset(soft bool) {
	s.normal.Reset(soft)
}
