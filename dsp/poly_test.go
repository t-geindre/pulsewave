package dsp

import (
	"testing"
	"time"
)

var voiceFact = func() *Voice {
	const sr = 44100.0

	freq := NewParam(440)
	env := NewADSR(sr, time.Millisecond*10, time.Millisecond*50, 0.8, time.Millisecond*100)
	src := NewOscillator(sr, ShapeNoise, freq, nil, nil) // Noise = fastest

	return NewVoice(src, freq, env)
}

func BenchmarkPoly_Process(b *testing.B) {
	poly := NewPolyVoice(16, voiceFact)
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
	poly := NewPolyVoice(16, voiceFact)
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
