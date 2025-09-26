package oscillator

type Merger struct {
	oscs  []Oscillator
	gains []float64
}

func NewMerger() *Merger {
	return &Merger{}
}

func (m *Merger) NextSample() float64 {
	v := 0.0

	for i, osc := range m.oscs {
		v += osc.NextSample() * m.gains[i]
	}

	return v
}

func (m *Merger) Append(osc Oscillator, gain float64) {
	m.oscs = append(m.oscs, osc)
	m.gains = append(m.gains, gain)
}

func (m *Merger) SetFreq(freq float64) {
	for _, osc := range m.oscs {
		osc.SetFreq(freq)
	}
}
