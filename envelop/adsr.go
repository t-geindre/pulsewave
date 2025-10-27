package envelop

import "time"

type EnvState int

const (
	EnvIdle EnvState = iota
	EnvAttack
	EnvDecay
	EnvSustain
	EnvRelease
)

type ADSR struct {
	attack  float64
	decay   float64
	sustain float64
	release float64

	sr    float64
	state EnvState
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
		state:   EnvIdle,
		value:   0,
		aStep:   aStep,
		dStep:   dStep,
		rStep:   0,
	}
}

func (e *ADSR) NoteOn(_, _ float64) {
	// Si attack <= 0, on "claque" au pic puis Decay
	if e.attack <= 0 {
		e.value = 1.0
		e.state = EnvDecay
		// dStep dépend du sustain (au cas où il a changé)
		if e.decay > 0 {
			e.dStep = (1.0 - e.sustain) / (e.decay * e.sr)
		} else {
			e.dStep = 0
		}
		return
	}

	// Soft-retrigger: on repart de la valeur courante sans discontinuité
	if e.value < 0 {
		e.value = 0
	}
	if e.value >= 1.0 {
		// Si on est déjà au max, on enchaîne Decay
		e.value = 1.0
		e.state = EnvDecay
	} else {
		e.state = EnvAttack
		// Temps d'attaque restant = (1 - value) * attack
		e.aStep = (1.0 - e.value) / (e.attack * e.sr)
	}

	// Toujours garder dStep synchro (au cas où sustain/decay ont changé)
	if e.decay > 0 {
		e.dStep = (1.0 - e.sustain) / (e.decay * e.sr)
	} else {
		e.dStep = 0
	}
}

func (e *ADSR) NoteOff(float64) {
	if e.release <= 0 || e.value <= 0 {
		e.value = 0
		e.state = EnvIdle
		e.rStep = 0
		return
	}
	e.state = EnvRelease
	// Release linéaire depuis le niveau courant
	e.rStep = e.value / (e.release * e.sr)
}

func (e *ADSR) NextValue() (float64, float64) {
	switch e.state {
	case EnvIdle:
		e.value = 0

	case EnvAttack:
		// aStep dépend du point de départ au moment du NoteOn
		e.value += e.aStep
		if e.value >= 1.0 || e.attack == 0 {
			e.value = 1.0
			e.state = EnvDecay
			// dStep peut avoir été recalculé dans NoteOn ; sinon on le sécurise
			if e.decay > 0 {
				e.dStep = (1.0 - e.sustain) / (e.decay * e.sr)
			} else {
				e.dStep = 0
			}
		}

	case EnvDecay:
		if e.decay == 0 {
			e.value = e.sustain
			e.state = EnvSustain
		} else {
			e.value -= e.dStep
			if e.value <= e.sustain {
				e.value = e.sustain
				e.state = EnvSustain
			}
		}

	case EnvSustain:
		// niveau constant

	case EnvRelease:
		e.value -= e.rStep
		if e.value <= 0 {
			e.value = 0
			e.state = EnvIdle
		}
	}
	return e.value, e.value
}

func (e *ADSR) Value() float64 { return e.value }

func (e *ADSR) IsActive() bool { return e.state != EnvIdle }

func (e *ADSR) Reset() {
	e.state = EnvIdle
	e.value = 0
}

func (e *ADSR) GetState() EnvState {
	return e.state
}
