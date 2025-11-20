package dsp

import (
	"testing"
)

var voiceFact = func() *Voice {
	const sr = 44100.0

	freq := NewParam(440)
	env := NewADSR(sr,
		NewConstParam(10/1000),
		NewConstParam(50/1000),
		NewConstParam(0.8),
		NewConstParam(100/1000),
	)
	src := NewNoise()

	return NewVoice(src, freq, env)
}

func BenchmarkPoly_Process(b *testing.B) {
	poly := NewPolyVoice(16, NewConstParam(16), NewConstParam(PolyStealOldest), voiceFact)
	for i := 0; i < 16; i++ {
		poly.NoteOn(i, 1.0)
	}

	var block Block
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		poly.Process(&block)
		block.Cycle++
	}
}

func TestPoly_ProcessNoAlloc(t *testing.T) {
	poly := NewPolyVoice(16, NewConstParam(16), NewConstParam(PolyStealOldest), voiceFact)
	for i := 0; i < 16; i++ {
		poly.NoteOn(i, 1.0)
	}

	var block Block

	allocs := testing.AllocsPerRun(100, func() {
		poly.Process(&block)
		block.Cycle++
	})

	if allocs != 0 {
		t.Errorf("expected 0 allocs, got %f", allocs)
	}
}
