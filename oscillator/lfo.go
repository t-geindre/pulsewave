package oscillator

type Lfo struct {
	Oscillator
	bind func(float64)
}

func NewLfo(osc Oscillator, bind func(float64)) *Lfo {
	return &Lfo{
		Oscillator: osc,
		bind:       bind,
	}
}

func (l *Lfo) NextSample() float64 {
	v := l.Oscillator.NextSample()
	l.bind(v)
	return 0
}
