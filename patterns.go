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

func CrazyFrogLeadPattern() *sequencer.Pattern {
	p := sequencer.NewPattern()

	p.Append(sequencer.D4, 4, 1, .85)
	p.Append(sequencer.F4, 4, 1, .85)
	p.Append(sequencer.D4, 2, 1, .85)
	p.Append(sequencer.D4, 1, 1, .85)
	p.Append(sequencer.G4, 2, 1, .85)
	p.Append(sequencer.D4, 1, 1, .85)
	p.Append(sequencer.C4, 2, 1, .85)

	p.Append(sequencer.D4, 4, 1, .85)
	p.Append(sequencer.A4, 4, 1, .85)
	p.Append(sequencer.D4, 2, 1, .85)
	p.Append(sequencer.D4, 1, 1, .85)
	p.Append(sequencer.As4, 2, 1, .85)
	p.Append(sequencer.A4, 1, 1, .85)
	p.Append(sequencer.F4, 2, 1, .85)

	p.Append(sequencer.D4, 2, 1, .85)
	p.Append(sequencer.A4, 2, 1, .85)
	p.Append(sequencer.D5, 2, 1, .85)
	p.Append(sequencer.D4, 1, 1, .85)
	p.Append(sequencer.C4, 2, 1, .85)
	p.Append(sequencer.C4, 1, 1, .85)
	p.Append(sequencer.A3, 2, 1, .85)

	p.Append(sequencer.E4, 2, 1, .85)
	p.Append(sequencer.D4, 8, 1, .85)

	return p
}

func CrazyFrogKickPattern() *sequencer.Pattern {
	p := sequencer.NewPattern()
	for i := 0; i < 14; i++ {
		p.Append(sequencer.C4, 4, 1, .95)
	}
	return p
}

func CrazyFrogHighHatPattern() *sequencer.Pattern {
	p := CrazyFrogKickPattern().Clone()
	p.Prepend(0, 2, 0, 0)
	p.Notes[len(p.Notes)-1].Length -= 2
	return p
}

func CrazyFrogBassPattern() *sequencer.Pattern {
	p := sequencer.NewPattern()

	rep := func(n int, f float64) {
		for i := 0; i < n; i++ {
			p.Append(f, 1, 1, .85)
		}
	}

	rep(4, sequencer.D4)
	rep(4, sequencer.F4)
	rep(4, sequencer.D4)
	rep(4, sequencer.C4)

	rep(4, sequencer.D4)
	rep(4, sequencer.A4)
	rep(4, sequencer.D4)
	rep(4, sequencer.C4)

	rep(2, sequencer.D4)
	rep(2, sequencer.A4)
	rep(2, sequencer.D5)
	rep(1, sequencer.D4)
	rep(2, sequencer.C4)
	rep(1, sequencer.C4)
	rep(2, sequencer.A3)
	rep(2, sequencer.E4)
	rep(8, sequencer.D4)

	return p
}
