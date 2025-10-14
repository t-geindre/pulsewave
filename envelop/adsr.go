package envelop

import "time"

type envState int

const (
	envIdle envState = iota
	envAttack
	envDecay
	envSustain
	envRelease
)

type ADSR struct {
	attack  float64
	decay   float64
	sustain float64
	release float64

	sr    float64
	state envState
	value float64
	aStep float64
	dStep float64
	rStep float64
}

func NewADSR(sampleRate float64, attack, decay, release time.Duration, sustain float64) *ADSR {
	if sustain < 0 {
		sustain = 0
	} else if sustain > 1 {
		sustain = 1
	}
	aStep := 0.0
	dStep := 0.0
	if attack > 0 {
		aStep = 1.0 / (attack.Seconds() * sampleRate)
	}
	if decay > 0 {
		dStep = (1.0 - sustain) / (decay.Seconds() * sampleRate)
	}
	return &ADSR{
		attack:  attack.Seconds(),
		decay:   decay.Seconds(),
		sustain: sustain,
		release: release.Seconds(),
		sr:      sampleRate,
		state:   envIdle,
		value:   0,
		aStep:   aStep,
		dStep:   dStep,
		rStep:   0,
	}
}

func (e *ADSR) NoteOn() {
	if e.attack <= 0 {
		e.value = 1
		e.state = envDecay
	} else {
		e.state = envAttack
	}
}

func (e *ADSR) NoteOff() {
	if e.release <= 0 {
		e.value = 0
		e.state = envIdle
		return
	}
	if e.value <= 0 {
		e.value = 0
		e.state = envIdle
		e.rStep = 0
		return
	}
	e.state = envRelease
	e.rStep = e.value / (e.release * e.sr) // linéaire jusqu'à 0
}

func (e *ADSR) Next() float64 {
	switch e.state {
	case envIdle:
		e.value = 0
	case envAttack:
		e.value += e.aStep
		if e.value >= 1.0 || e.attack == 0 {
			e.value = 1.0
			e.state = envDecay
		}
	case envDecay:
		if e.decay == 0 {
			e.value = e.sustain
			e.state = envSustain
		} else {
			e.value -= e.dStep
			if e.value <= e.sustain {
				e.value = e.sustain
				e.state = envSustain
			}
		}
	case envSustain:
		// Keep sustain level
	case envRelease:
		e.value -= e.rStep
		if e.value <= 0 {
			e.value = 0
			e.state = envIdle
		}
	}
	return e.value
}

func (e *ADSR) Value() float64 { return e.value }

func (e *ADSR) IsActive() bool { return e.state != envIdle }
