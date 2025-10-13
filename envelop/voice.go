package envelop

import (
	"synth/oscillator"
)

type Voice struct {
	osc  oscillator.Oscillator
	env  *ADSR
	gain float64
}

func NewVoice(sr float64, osc oscillator.Oscillator, adsr *ADSR) *Voice {
	return &Voice{
		osc:  osc,
		env:  adsr,
		gain: 1,
	}
}

func (v *Voice) NoteOn(freq, velocity float64) {
	v.osc.SetFreq(freq)
	v.osc.ResetPhase()
	if velocity <= 0 {
		velocity = 1
	}
	if velocity > 1 {
		velocity = 1
	}
	v.gain = velocity
	v.env.NoteOn()
}

func (v *Voice) NoteOff() {
	v.env.NoteOff()
}

func (v *Voice) NextSample() float64 {
	amp := v.env.Next()
	s := v.osc.NextSample() * amp
	return s
}

func (v *Voice) IsActive() bool {
	return v.env.IsActive()
}
