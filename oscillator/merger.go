package oscillator

import "math"

type Merger struct {
	oscs       []Oscillator
	gains      []float64
	sumAbsGain float64
}

func NewMerger() *Merger {
	return &Merger{}
}

func (m *Merger) NextSample() float64 {
	v := 0.0

	for i, osc := range m.oscs {
		v += osc.NextSample() * m.gains[i]
	}

	return v / m.sumAbsGain
}

func (m *Merger) Append(osc Oscillator, gain float64) {
	m.oscs = append(m.oscs, osc)
	m.gains = append(m.gains, gain)

	s := 0.0
	for _, g := range m.gains {
		s += math.Abs(g)
	}
	m.sumAbsGain = s
}

func (m *Merger) SetFreq(freq float64) {
	for _, osc := range m.oscs {
		osc.SetFreq(freq)
	}
}

func (m *Merger) ResetPhase() {
	for _, osc := range m.oscs {
		osc.ResetPhase()
	}
}

func (m *Merger) SetPhaseShift(phase float64) {
	for _, osc := range m.oscs {
		osc.SetPhaseShift(phase)
	}
}
