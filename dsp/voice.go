package dsp

import "synth/audio"

type Envelope interface {
	NoteOn()
	NoteOff()
	IsIdle() bool
	ParamModulator
}

type Voice struct {
	audio.Source
	freq Param
	env  Envelope
}

func NewVoice(src audio.Source, freq Param, env Envelope) *Voice {
	gain := NewParam(0)
	*gain.ModInputs() = append(*gain.ModInputs(), NewModInput(env, 1.0, nil))

	return &Voice{
		Source: NewVca(src, gain),
		freq:   freq,
		env:    env,
	}
}

func (v *Voice) NoteOne(freq, vel float32) {
	v.env.NoteOn()
}

func (v *Voice) NoteOff() {
	v.env.NoteOff()
}
