package envelop

import (
	"synth/audio"
)

type Voice struct {
	src  audio.Source
	env  audio.Source
	gain float64
}

func NewVoice(src, env audio.Source) *Voice {
	return &Voice{
		src:  src,
		env:  env,
		gain: 1,
	}
}

func (v *Voice) NoteOn(freq, velocity float64) {
	v.src.NoteOn(freq, velocity)
	v.env.NoteOn(freq, velocity)
	v.gain = velocity
}

func (v *Voice) NoteOff() {
	v.src.NoteOff()
	v.env.NoteOff()
}

func (v *Voice) NextValue() (float64, float64) {
	amp, _ := v.env.NextValue()
	l, r := v.src.NextValue()
	return l * amp * v.gain, r * amp * v.gain
}

func (v *Voice) IsActive() bool {
	return v.env.IsActive()
}

func (v *Voice) Reset() {
	v.src.Reset()
	v.env.Reset()
}
