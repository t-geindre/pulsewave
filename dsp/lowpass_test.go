package dsp

import "testing"

func TestLowPassSVF_ProcessNoAlloc(t *testing.T) {
	const sr = 44100.0
	src := NewNoise()
	filter := NewLowPassSVF(sr, src, NewConstParam(1000), NewConstParam(0.707))

	var block Block

	allocs := testing.AllocsPerRun(1000, func() {
		filter.Process(&block)
		block.Cycle++
	})

	if allocs != 0 {
		t.Errorf("expected 0 allocations, got %v", allocs)
	}
}
