package dsp

import (
	"testing"
	"time"
)

func BenchmarkPoly_Process(b *testing.B) {
	const sr = 44100.0

	voiceFact := func() *Voice {
		freq := NewParam(440)
		env := NewADSR(sr, time.Millisecond*10, time.Millisecond*50, 0.8, time.Millisecond*100)
		src := NewOscillator(sr, ShapeNoise, freq, nil, nil) // Noise = fastest
		return NewVoice(src, freq, env)
	}

	poly := NewPolyVoice(16, voiceFact)
	for i := 0; i < 16; i++ {
		poly.NoteOn(i, 1.0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var block Block
		poly.Process(&block)
	}
}
