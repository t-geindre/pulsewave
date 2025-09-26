package main

import (
	"math/rand"
	"synth/audio"
	"synth/automation"
	"synth/envelop"
	"synth/oscillator"
	"time"
)

const SampleRate = 44100

func main() {
	osc := oscillator.NewMerger()

	oscBase := oscillator.NewSine(SampleRate, 440)
	osc.Append(oscBase, 1)

	oscBass := oscillator.NewTuned(oscillator.NewSine(SampleRate, 220))
	oscBass.SetOctaveOffset(-1)
	osc.Append(oscBass, 1)

	adsr := envelop.NewADSR(SampleRate, time.Millisecond*20, time.Millisecond*50, time.Millisecond*50, 1)
	voice := envelop.NewVoice(SampleRate, osc, adsr)

	cMajorScale := []float64{261.63, 293.66, 329.63, 349.23, 392.00, 440.00, 493.88, 523.25}

	var noteOn, noteOff automation.NextFunc
	noteOn = func(t time.Duration) (automation.NextFunc, time.Duration) {
		voice.NoteOn(cMajorScale[rand.Intn(len(cMajorScale))], 1)
		return noteOff, t + time.Millisecond*50
	}
	noteOff = func(t time.Duration) (automation.NextFunc, time.Duration) {
		voice.NoteOff()
		return noteOn, t + time.Millisecond*100
	}
	automate := automation.NewTimed(SampleRate, noteOn)

	tracks := audio.NewTrackSet(voice, automate)
	player := audio.NewPlayer(SampleRate, tracks)

	for player.IsPlaying() {
		time.Sleep(time.Millisecond * 100)
	}
}
