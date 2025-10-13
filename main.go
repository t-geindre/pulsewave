package main

import (
	"synth/audio"
	"synth/effect"
	"synth/envelop"
	"synth/oscillator"
	"synth/sequencer"
	"time"
)

const SampleRate = 44100

func main() {
	tracks := audio.NewTrackSet()

	lead := sequencer.NewSequencer(SampleRate, 140, 8, 4, newLeadVoice)
	lead.SetPattern(RandomCMajorPattern(64, []int{2, 4}, .75))
	lead.SetLoopMode(sequencer.LoopSoft)

	delay := effect.NewFeedbackDelay(SampleRate, lead)
	delay.SetDelay(lead.GetBeatDuration() / 2)
	delay.SetMix(0.3)
	delay.SetFeedback(0.5)

	lp := effect.NewLowPassFilter(SampleRate, delay)
	lp.SetCutoffHz(1000)
	lp.SetQ(5)

	tracks.Append(lp)

	// LFO
	lfo := oscillator.NewLfo(oscillator.NewSine(SampleRate, .1), func(v float64) {
		cutoff := 500 + (v+1)*0.5*5000 // 500Hz - 5500Hz
		lp.SetCutoffHz(cutoff)
	})
	tracks.Append(lfo)
	// --

	player := audio.NewPlayer(SampleRate, tracks)

	for player.IsPlaying() {
		time.Sleep(100 * time.Millisecond)
	}
}

func newLeadVoice() sequencer.Voice {
	merged := oscillator.NewMerger()

	for i := 0; i < 3; i++ {
		osc := oscillator.NewSaw(SampleRate, 110.0)
		osc.SetPhaseShift(float64(i) * 0.01)
		tuned := oscillator.NewTuned(osc)
		tuned.SetDetuneCents(float64(i) * 8.0)
		merged.Append(tuned, 1.0)
	}

	adsr := envelop.NewADSR(SampleRate, time.Millisecond*25, time.Millisecond*100, time.Millisecond*50, .8)
	return envelop.NewVoice(SampleRate, merged, adsr)
}
