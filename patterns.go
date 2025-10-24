package main

import (
	"math/rand"
	"synth/sequencer"
)

func RandomCMajorPattern(want int, lengths []int, gate float64) *sequencer.Pattern {
	freqs := []float64{
		sequencer.C4,
		sequencer.F4,
		sequencer.G4,
		sequencer.C5,
	}

	length := 0
	pattern := sequencer.NewPattern()
	for length < want {
		l := 0
		for {
			l = lengths[rand.Intn(len(lengths))]
			if length+l <= want {
				break
			}
		}
		length += l
		f := freqs[rand.Intn(len(freqs))]
		pattern.Append(f, l, 1, gate)
	}
	return pattern
}

func CMajorScalePattern() *sequencer.Pattern {
	const L = 4
	const G = .75

	pattern := sequencer.NewPattern()
	pattern.Append(sequencer.C4, L, 1, G)
	pattern.Append(sequencer.D4, L, 1, G)
	pattern.Append(sequencer.E4, L, 1, G)
	pattern.Append(sequencer.F4, L, 1, G)
	pattern.Append(sequencer.G4, L, 1, G)
	pattern.Append(sequencer.A4, L, 1, G)
	pattern.Append(sequencer.B4, L, 1, G)
	pattern.Append(sequencer.C5, L, 1, G)
	return pattern
}

func TetrisThemeAPattern() *sequencer.Pattern {
	p := sequencer.NewPattern()

	p.Append(sequencer.E5, 4, 1, .75)
	p.Append(sequencer.B4, 2, 1, .75)
	p.Append(sequencer.C5, 2, 1, .75)
	p.Append(sequencer.D5, 4, 1, .75)
	p.Append(sequencer.C5, 2, 1, .75)
	p.Append(sequencer.B4, 2, 1, .75)
	p.Append(sequencer.A4, 4, 1, .75)
	p.Append(sequencer.A4, 2, 1, .75)
	p.Append(sequencer.C5, 2, 1, .75)
	p.Append(sequencer.E5, 4, 1, .75)
	p.Append(sequencer.D5, 2, 1, .75)
	p.Append(sequencer.C5, 2, 1, .75)
	p.Append(sequencer.B4, 4, 1, .75)
	p.Append(sequencer.B4, 2, 1, .75)
	p.Append(sequencer.C5, 2, 1, .75)
	p.Append(sequencer.D5, 4, 1, .75)
	p.Append(sequencer.E5, 4, 1, .75)
	p.Append(sequencer.C5, 4, 1, .75)
	p.Append(sequencer.A4, 4, 1, .75)
	p.Append(sequencer.A4, 4, 1, .75)

	return p
}

func PiratesBassPattern() *sequencer.Pattern {
	p := sequencer.NewPattern()

	p.Append(0, 2, 1, .75)
	p.Append(sequencer.D3, 4, 1, .75)

	p.Append(0, 2, 1, .75)
	p.Append(sequencer.F3, 4, 1, .75)

	p.Append(sequencer.F4, 1, 1, .75)
	p.Append(sequencer.G4, 1, 1, .75)
	p.Append(sequencer.E4, 2, 1, .75)
	p.Append(sequencer.E4, 2, 1, .75)

	p.Append(sequencer.D4, 1, 1, .75)
	p.Append(sequencer.C4, 1, 1, .75)
	p.Append(sequencer.D4, 2, 1, .75)

	return p
}
