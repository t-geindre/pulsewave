package preset

import (
	"synth/dsp"
	"testing"
)

func TestPolysynth_ProcessNoAlloc(t *testing.T) {
	synth := NewPolysynth(44100)

	for i := 0; i < 16; i++ {
		synth.voice.NoteOn(10+i, 1.0)
	}
	var block dsp.Block

	allocs := testing.AllocsPerRun(1000, func() {
		synth.Process(&block)
	})
	if allocs != 0 {
		t.Errorf("expected 0 allocations, got %v", allocs)
	}
}
