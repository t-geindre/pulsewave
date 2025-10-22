package song

import (
	"synth/audio"
	"synth/effect"
	"synth/envelop"
	"synth/oscillator"
	"synth/sequencer"
	"time"
)

type CrazyFrog struct {
	*audio.TrackSet
	sr float64
}

func NewCrazyFrog(SampleRate float64) audio.Source {
	const BPM = 120
	const SPB = 4 // Step per beat
	const Repeat = 10

	c := &CrazyFrog{
		TrackSet: audio.NewTrackSet(),
		sr:       SampleRate,
	}

	// Lead
	lead := sequencer.NewSequencer(SampleRate, BPM, 4, SPB, c.leadVoiceFactory)
	lead.AppendAndRepeat(c.leadPattern(), Repeat)

	leadDelay := effect.NewFeedback(SampleRate, lead)
	leadDelay.SetDelay(lead.GetBeatDuration() / 2)
	leadDelay.SetMix(.4)
	leadDelay.SetFeedback(.3)

	c.Append(leadDelay, .8)

	// Bass
	bass := sequencer.NewSequencer(SampleRate, BPM, 4, SPB, c.bassVoiceFactory)
	bass.AppendAndRepeat(c.bassPattern(), Repeat)

	lpfBass := effect.NewLowPassFilter(SampleRate, bass)
	lpfBass.SetCutoffHz(800)
	lpfBass.SetQ(.7)

	c.Append(lpfBass, .08)

	// High hat
	hh := sequencer.NewSequencer(SampleRate, BPM, 4, SPB, c.highHatVoiceFactory)
	hh.AppendAndRepeat(c.highHatPattern(), 10)
	
	c.Append(hh, .2)

	// Kick
	kick := sequencer.NewSequencer(SampleRate, BPM, 4, SPB, c.kickVoiceFactory)
	kick.AppendAndRepeat(c.kickPattern(), Repeat)
	c.Append(kick, .9)

	return c
}

func (c *CrazyFrog) leadVoiceFactory() audio.Source {
	// Saw unison
	uni := effect.NewUnison(func() audio.Source {
		osc := oscillator.NewSaw(c.sr)
		return osc
	}, 8, 12, 0.9, .1, 0.75)

	// Bass line on lead
	bass := effect.NewTuner(oscillator.NewTriangle(c.sr))
	bass.SetOctaveOffset(-1)

	merger := effect.NewMerger()
	merger.Append(uni, 1, 1)
	merger.Append(bass, .5, .5)

	adsr := envelop.NewADSR(c.sr, time.Millisecond*5, time.Millisecond*100, time.Millisecond*50, .8)
	voice := envelop.NewVoice(merger, adsr)

	return voice
}

func (c *CrazyFrog) bassVoiceFactory() audio.Source {
	// Unison saw
	uni := effect.NewUnison(func() audio.Source {
		osc := oscillator.NewSquare(c.sr)
		return osc
	}, 6, 8, 0.9, .9, 0.75)

	// Octave down
	tuner := effect.NewTuner(uni)
	tuner.SetOctaveOffset(-2)

	adsr := envelop.NewADSR(c.sr, time.Millisecond*10, time.Millisecond*10, time.Millisecond*10, .7)
	return envelop.NewVoice(tuner, adsr)
}

func (c *CrazyFrog) highHatVoiceFactory() audio.Source {
	noise := oscillator.NewNoise()
	adsr := envelop.NewADSR(c.sr, time.Millisecond*5, time.Millisecond*70, 0, 0)
	return envelop.NewVoice(noise, adsr)
}

func (c *CrazyFrog) kickVoiceFactory() audio.Source {
	merger := effect.NewMerger()

	picthMod := envelop.NewADSR(c.sr, 0, time.Millisecond*100, 0, 0.0)
	merger.Append(picthMod, 0, 0)

	sine := oscillator.NewSine(c.sr)
	kick := audio.NewCallbackSrc(func() (float64, float64) {
		v, _ := picthMod.NextValue()
		sine.SetFreq(60.0 + v*300.0)
		return sine.NextValue()
	})
	merger.Append(kick, 1, 1)

	adsr := envelop.NewADSR(c.sr, time.Millisecond*1, time.Millisecond*450, 0, 0.0)
	return envelop.NewVoice(merger, adsr)
}

func (_ *CrazyFrog) leadPattern() *sequencer.Pattern {
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
	p.Append(sequencer.D4, 10, 1, .85/4*3)

	return p
}

func (_ *CrazyFrog) kickPattern() *sequencer.Pattern {
	p := sequencer.NewPattern()
	for i := 0; i < 14; i++ {
		p.Append(sequencer.C4, 4, 1, .95)
	}
	return p
}

func (c *CrazyFrog) highHatPattern() *sequencer.Pattern {
	p := c.kickPattern().Clone()
	p.Prepend(0, 2, 0, 0)
	return p
}

func (_ *CrazyFrog) bassPattern() *sequencer.Pattern {
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
	rep(2, sequencer.D4)
	rep(2, sequencer.D4)
	rep(1, sequencer.D4)
	rep(2, sequencer.C4)
	rep(1, sequencer.C4)
	rep(2, sequencer.A4)
	rep(2, sequencer.E4)
	rep(10, sequencer.D4)

	return p
}
