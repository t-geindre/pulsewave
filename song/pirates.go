package song

import (
	"synth/audio"
	"synth/effect"
	"synth/envelop"
	"synth/oscillator"
	"synth/sequencer"
	"time"
)

type Pirates struct {
	*audio.TrackSet
	sr float64
}

func NewPirates(SampleRate float64) audio.Source {
	const BPM = 120

	p := &Pirates{
		TrackSet: audio.NewTrackSet(),
		sr:       SampleRate,
	}

	// Lead
	seq := sequencer.NewSequencer(SampleRate, BPM, 4, 6, p.leadVoice)

	seq.Append(p.leadPatternA())
	seq.Append(p.leadPatternA())
	seq.Append(p.leadPatternB())

	seq.Append(p.leadPatternA())
	seq.Append(p.leadPatternA())
	seq.Append(p.leadPatternB())

	delay := effect.NewFeedback(SampleRate, seq)
	delay.SetDelay(time.Millisecond * 500)
	delay.SetFeedback(.2)
	delay.SetMix(.3)

	p.Append(delay, 1)

	// Bass
	bassSeq := sequencer.NewSequencer(SampleRate, BPM, 4, 6, p.bassVoice)
	bassSeq.Append(p.bassPattern(true))
	bassSeq.Append(p.bassPattern(false))

	p.Append(bassSeq, .6)

	return p
}

func (p *Pirates) leadVoice() audio.Source {
	// Saw + unison
	uni := effect.NewUnison(func() audio.Source {
		return oscillator.NewSaw(p.sr)
	}, 8, 8, 0.9, 0, .75)

	// LPF + ADSR mod
	merge := effect.NewMerger()

	cutoffMod := envelop.NewADSR(p.sr, 0, time.Millisecond*250, 0, 0)
	merge.Append(cutoffMod, 0, 0)

	lpf := effect.NewLowPassFilter(p.sr, uni)
	lpf.SetQ(.7)
	merge.Append(lpf, 1, 1)

	cutoff := audio.NewCallbackSrc(func() (float64, float64) {
		v, _ := cutoffMod.NextValue()
		lpf.SetCutoffHz(800 + v*4000.0)
		return 0, 0
	})
	merge.Append(cutoff, 0, 0)

	return envelop.NewVoice(
		merge,
		envelop.NewADSR(p.sr, time.Millisecond*10, time.Millisecond*150, time.Millisecond*250, 1),
	)
}

func (p *Pirates) bassVoice() audio.Source {
	uni := effect.NewUnison(func() audio.Source {
		osc := oscillator.NewTriangle(p.sr)
		return osc
	}, 8, 24, 0.9, 0.9, .75)

	return envelop.NewVoice(
		uni,
		envelop.NewADSR(p.sr, time.Millisecond*5, time.Millisecond*100, time.Millisecond*100, 1),
	)
}

func (_ *Pirates) leadPatternA() *sequencer.Pattern {
	p := sequencer.NewPattern()

	p.Append(sequencer.A3, 2, 1, .75)
	p.Append(sequencer.C4, 2, 1, .75)

	p.Append(sequencer.D4, 4, 1, .75)
	p.Append(sequencer.D4, 4, 1, .75)
	p.Append(sequencer.D4, 2, 1, .75)
	p.Append(sequencer.E4, 2, 1, .75)

	p.Append(sequencer.F4, 4, 1, .75)
	p.Append(sequencer.F4, 4, 1, .75)
	p.Append(sequencer.F4, 2, 1, .75)
	p.Append(sequencer.G4, 2, 1, .75)

	p.Append(sequencer.E4, 4, 1, .75)
	p.Append(sequencer.E4, 4, 1, .75)
	p.Append(sequencer.D4, 2, 1, .75)
	p.Append(sequencer.C4, 2, 1, .75)

	p.Append(sequencer.D4, 8, 1, .75) // length +2*2 at beginning = 12

	return p
}

func (_ *Pirates) leadPatternB() *sequencer.Pattern {
	p := sequencer.NewPattern()

	p.Append(sequencer.A3, 2, 1, .75)
	p.Append(sequencer.C4, 2, 1, .75)

	p.Append(sequencer.D4, 4, 1, .75)
	p.Append(sequencer.D4, 4, 1, .75)
	p.Append(sequencer.D4, 2, 1, .75)
	p.Append(sequencer.F4, 2, 1, .75)

	p.Append(sequencer.G4, 4, 1, .75)
	p.Append(sequencer.G4, 4, 1, .75)
	p.Append(sequencer.G4, 2, 1, .75)
	p.Append(sequencer.A4, 2, 1, .75)

	p.Append(sequencer.As4, 4, 1, .75)
	p.Append(sequencer.As4, 4, 1, .75)
	p.Append(sequencer.A4, 2, 1, .75)
	p.Append(sequencer.G4, 2, 1, .75)

	p.Append(sequencer.A4, 2, 1, .75)
	p.Append(sequencer.D4, 6, 1, .75)
	p.Append(sequencer.D4, 2, 1, .75)
	p.Append(sequencer.E4, 2, 1, .75)

	p.Append(sequencer.F4, 4, 1, .75)
	p.Append(sequencer.F4, 4, 1, .75)
	p.Append(sequencer.G4, 4, 1, .75)

	p.Append(sequencer.A4, 2, 1, .75)
	p.Append(sequencer.D4, 6, 1, .75)
	p.Append(sequencer.D4, 2, 1, .75)
	p.Append(sequencer.F4, 2, 1, .75)

	p.Append(sequencer.E4, 4, 1, .75)
	p.Append(sequencer.E4, 4, 1, .75)
	p.Append(sequencer.F4, 2, 1, .75)
	p.Append(sequencer.D4, 2, 1, .75)

	p.Append(sequencer.E4, 8, 1, .75) // length +2*2 at beginning = 12

	return p
}

func (_ *Pirates) bassPattern(intro bool) *sequencer.Pattern {
	p := sequencer.NewPattern()

	if intro {
		p.Append(0, 4, 1, .75)
	}

	p.Append(sequencer.D3, 12, 1, 1)
	p.Append(sequencer.As2, 12, 1, 1)
	p.Append(sequencer.A2, 12, 1, 1)
	p.Append(sequencer.D3, 12, 1, 1)

	p.Append(sequencer.As2, 12, 1, 1)
	p.Append(sequencer.F2, 12, 1, 1)
	p.Append(sequencer.C3, 12, 1, 1)
	p.Append(sequencer.D3, 12, 1, 1)

	p.Append(sequencer.D3, 12, 1, 1)
	p.Append(sequencer.As2, 12, 1, 1)
	p.Append(sequencer.G2, 12, 1, 1)
	p.Append(sequencer.F2, 12, 1, 1)
	p.Append(sequencer.As2, 12, 1, 1)
	p.Append(sequencer.F2, 12, 1, 1)
	p.Append(sequencer.G2, 12, 1, 1)
	p.Append(sequencer.A2, 12, 1, 1)

	return p
}
