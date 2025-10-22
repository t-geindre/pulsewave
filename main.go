package main

import (
	"synth/audio"
	"synth/song"
	"time"
)

const SampleRate = 44100
const BPM = 120

func main() {
	player := audio.NewPlayer(SampleRate, song.NewCrazyFrog(SampleRate))

	for player.IsPlaying() {
		time.Sleep(100 * time.Millisecond)
	}
}

/*


func main() {
	tracks := audio.NewTrackSet()
	tracks.SetLoop(true)

	lead := sequencer.NewSequencer(SampleRate, BPM, 4, 4, newLeadVoice)
	lead.Append(CrazyFrogLeadPattern())

	leadDelay := effect.NewFeedback(SampleRate, lead)
	leadDelay.SetDelay(lead.GetBeatDuration() / 2)
	leadDelay.SetMix(.3)
	leadDelay.SetFeedback(.4)
	tracks.Append(leadDelay, 1)

	kicks := sequencer.NewSequencer(SampleRate, BPM, 2, 4, newKickVoice)
	kicks.Append(CrazyFrogKickPattern())
	tracks.Append(kicks, 1)

	highHat := sequencer.NewSequencer(SampleRate, BPM, 1, 4, newHighHatVoice)
	highHat.Append(CrazyFrogHighHatPattern())
	tracks.Append(highHat, .2)

	bass := sequencer.NewSequencer(SampleRate, BPM, 4, 4, newBassVoice)
	bass.Append(CrazyFrogBassPattern())
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
		tuned := oscillator.NewTuner(osc)
		tuned.SetDetuneCents(float64(i) * 8.0)
		merged.Append(tuned, 1.0)
	}

	sq := oscillator.NewSquare(SampleRate, 220.0)
	sq.SetPulseWidth(0.25)
	tq := oscillator.NewTuner(sq)
	merged.Append(tq, 0.5)

	tr := oscillator.NewSine(SampleRate, 110.0)
	tt := oscillator.NewTuner(tr)
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



func newHighHatVoice() sequencer.Voice {
	noise := oscillator.NewNoise()
	adsr := envelop.NewADSR(SampleRate, time.Millisecond*5, time.Millisecond*70, 0, 0)
	return envelop.NewVoice(SampleRate, noise, adsr)
}
*/
