package envelop

import (
	"synth/audio"
)

type Voice struct {
	src  audio.Source
	envs []audio.Source
	gain float64
}

func NewVoice(src audio.Source, envs ...audio.Source) *Voice {
	return &Voice{
		src:  src,
		envs: envs,
		gain: 1,
	}
}

func (v *Voice) NoteOn(freq, velocity float64) {
	v.src.NoteOn(freq, velocity)
	for _, e := range v.envs {
		e.NoteOn(freq, velocity)
	}
	v.gain = velocity // TODO linear softening to target gain
}

func (v *Voice) NoteOff(freq float64) {
	v.src.NoteOff(freq)
	for _, e := range v.envs {
		e.NoteOff(freq)
	}
}

func (v *Voice) NextValue() (float64, float64) {
	amp := 1.0
	if len(v.envs) > 0 {
		amp, _ = v.envs[0].NextValue()
	}
	l, r := v.src.NextValue()
	return l * amp * v.gain, r * amp * v.gain
}

func (v *Voice) IsActive() bool {
	for _, e := range v.envs {
		if e.IsActive() {
			return true
		}
	}
	return false
}

func (v *Voice) Reset() {
	v.src.Reset()
	for _, e := range v.envs {
		e.Reset()
	}
}
