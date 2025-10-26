package instrument

import (
	"synth/audio"
	"synth/effect"
	"synth/envelop"
	"synth/oscillator"
	"time"
)

type Synth struct {
	audio.Source
}

func NewSynth(SampleRate float64) *Synth {
	// Oscillator
	oscVoice := func() audio.Source {
		merger := effect.NewMerger()
		merger.Append(oscillator.NewSaw(SampleRate), 1, 1)
		merger.Append(oscillator.NewSquare(SampleRate), 1, 1)
		merger.Append(oscillator.NewSine(SampleRate), 1, 1)
		return merger
	}

	// Unison
	unison := effect.NewUnison(oscVoice, 8, 12, .9, 0, .75)

	// LPF + ADSR modulation
	lpf := effect.NewLowPassFilter(SampleRate, unison)
	lpf.SetQ(.7)

	lpfModSrc := oscillator.NewSine(SampleRate) //envelop.NewADSR(SampleRate, time.Millisecond*0, time.Millisecond*1500, time.Millisecond*1500, 0)
	lpfModSrc.SetFreq(.1)
	lpfModSrc.SetPhaseShift(.25)

	lpfModApp := audio.NewCallbackSrc(func() (L, R float64) {
		v, _ := lpfModSrc.NextValue()
		v = (v + 1) / 2 // Normalize to [0,1]
		lpf.SetCutoffHz(500 + v*4000)
		return 0, 0
	})

	merger := effect.NewMerger()
	merger.Append(lpf, 1, 1)
	merger.Append(lpfModApp, 0, 0)

	// Global envelope
	adsr := envelop.NewADSR(SampleRate, time.Millisecond*50, time.Millisecond*100, time.Millisecond*200, .9)

	// Voice
	voice := envelop.NewVoice(merger, adsr, lpfModSrc)

	// Delay
	delay := effect.NewFeedback(SampleRate, voice)
	delay.SetDelay(time.Millisecond * 300)
	delay.SetMix(.3)
	delay.SetFeedback(.4)

	return &Synth{
		Source: delay,
	}
}
