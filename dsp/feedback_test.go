package dsp

import "testing"

func TestFeedbackDelay_ProcessNoAlloc(t *testing.T) {
	const sr = 44100.0
	src := NewNoise()
	fb := NewFeedbackDelay(
		sr, 2.0, src,
		NewConstParam(0.5), NewConstParam(0.5), NewConstParam(0.5), NewConstParam(0.5),
	)

	var block Block

	allocs := testing.AllocsPerRun(1000, func() {
		fb.Process(&block)
		block.Cycle++
	})

	if allocs != 0 {
		t.Errorf("expected 0 allocations, got %v", allocs)
	}
}
