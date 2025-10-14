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
	tracks.SetLoop(true)

	lead := sequencer.NewSequencer(SampleRate, 120, 4, 4, newLeadVoice)
	lead.SetPattern(CrazyFrogLeadPattern())

	leadDelay := effect.NewFeedbackDelay(SampleRate, lead)
	leadDelay.SetDelay(lead.GetBeatDuration() / 2)
	leadDelay.SetMix(0.3)
	leadDelay.SetFeedback(0.4)
	tracks.Append(leadDelay, 1)

	kicks := sequencer.NewSequencer(SampleRate, 120, 2, 4, newKickVoice)
	kicks.SetPattern(CrazyFrogKickPattern())
	tracks.Append(kicks, .8)

	highHat := sequencer.NewSequencer(SampleRate, 120, 8, 4, newHighHatVoice)
	highHat.SetPattern(CrazyFrogHighHatPattern())
	tracks.Append(highHat, .2)

	bass := sequencer.NewSequencer(SampleRate, 120, 4, 4, newBassVoice)
	bass.SetPattern(CrazyFrogBassPattern())
	tracks.Append(bass, .4)

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

	sq := oscillator.NewSquare(SampleRate, 220.0)
	sq.SetPulseWidth(0.25)
	tq := oscillator.NewTuned(sq)
	merged.Append(tq, 0.5)

	tr := oscillator.NewSine(SampleRate, 110.0)
	tt := oscillator.NewTuned(tr)
	tt.SetOctaveOffset(-1)
	merged.Append(tt, 0.8)

	adsr := envelop.NewADSR(SampleRate, time.Millisecond*5, time.Millisecond*100, time.Millisecond*50, .8)
	return envelop.NewVoice(SampleRate, merged, adsr)
}

func newKickVoice() sequencer.Voice {
	sine := oscillator.NewSine(SampleRate, 20.0)

	pitchMod := envelop.NewADSR(SampleRate, 0, time.Millisecond*50, 0, 0.0)
	kick := oscillator.NewCallbackOscillator(func() float64 {
		sine.SetFreq(60.0 + pitchMod.Next()*300.0)
		return sine.NextSample()
	})

	adsr := envelop.NewADSR(SampleRate, time.Millisecond*1, time.Millisecond*500, 0, 0.0)
	return envelop.NewMultiEnvVoice(envelop.NewVoice(SampleRate, kick, adsr), pitchMod)
}

func newBassVoice() sequencer.Voice {
	merger := oscillator.NewMerger()
	for i := 0; i < 8; i++ {
		osc := oscillator.NewSaw(SampleRate, 110.0)
		tuned := oscillator.NewTuned(osc)
		tuned.SetDetuneCents(float64(i) * 2.0)
		tuned.SetOctaveOffset(-1)
		merger.Append(tuned, 1.0)
	}

	sine := oscillator.NewSine(SampleRate, 110.0)
	ts := oscillator.NewTuned(sine)
	ts.SetOctaveOffset(-2)
	merger.Append(ts, 0.4)

	adsr := envelop.NewADSR(SampleRate, time.Millisecond*10, time.Millisecond*10, time.Millisecond*10, .7)
	return envelop.NewVoice(SampleRate, merger, adsr)
}

func newHighHatVoice() sequencer.Voice {
	noise := oscillator.NewNoise()
	adsr := envelop.NewADSR(SampleRate, time.Millisecond*10, time.Millisecond*50, 0, 0)
	return envelop.NewVoice(SampleRate, noise, adsr)
}
