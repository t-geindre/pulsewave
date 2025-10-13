package main

import (
	"math/rand"
	"synth/sequencer"
)

func RandomCMajorPattern(want int, lengths []int, gate float64) *sequencer.Pattern {
	freqs := []float64{
		sequencer.C4,
		//sequencer.D4,
		//sequencer.E4,
		sequencer.F4,
		sequencer.G4,
		//sequencer.B4,
		//sequencer.A4,
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
	pattern := sequencer.NewPattern()
	pattern.Append(sequencer.C4, 4, 1, .85)
	pattern.Append(sequencer.D4, 4, 1, .85)
	pattern.Append(sequencer.E4, 4, 1, .85)
	pattern.Append(sequencer.F4, 4, 1, .85)
	pattern.Append(sequencer.G4, 4, 1, .85)
	pattern.Append(sequencer.A4, 4, 1, .85)
	pattern.Append(sequencer.B4, 4, 1, .85)
	pattern.Append(sequencer.C5, 4, 1, .85)

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

func ChildrenIntroLeadPattern() *sequencer.Pattern {
	p := sequencer.NewPattern()

	p.Append(sequencer.E4, 4, 1, .5)
	p.Append(sequencer.E4, 4, 1, .5)
	p.Append(sequencer.E4, 4, 1, .5)
	p.Append(sequencer.E4, 4, 1, .5)

	p.Append(sequencer.E4, 4, 1, .5)
	p.Append(sequencer.E4, 4, 1, .5)
	p.Append(sequencer.E4, 4, 1, .5)
	p.Append(sequencer.E4, 4, 1, .5)

	p.Append(sequencer.G4, 4, 1, .5)
	p.Append(sequencer.FS4, 2, 1, .5)

	p.Append(sequencer.D4, 4, 1, .5)
	p.Append(sequencer.D4, 4, 1, .5)
	p.Append(sequencer.D4, 4, 1, .5)
	p.Append(sequencer.D4, 4, 1, .5)

	p.Append(sequencer.D4, 4, 1, .5)
	p.Append(sequencer.D4, 4, 1, .5)
	p.Append(sequencer.D4, 4, 1, .5)
	p.Append(sequencer.D4, 4, 1, .5)

	p.Append(sequencer.G4, 4, 1, .5)
	p.Append(sequencer.FS4, 2, 1, .5)

	p.Append(sequencer.C4, 4, 1, .5)
	p.Append(sequencer.C4, 4, 1, .5)
	p.Append(sequencer.C4, 4, 1, .5)
	p.Append(sequencer.C4, 4, 1, .5)

	p.Append(sequencer.C4, 4, 1, .5)
	p.Append(sequencer.C4, 4, 1, .5)
	p.Append(sequencer.C4, 4, 1, .5)
	p.Append(sequencer.C4, 4, 1, .5)

	return p
}
