package oscillator

type CallbackOscillator struct {
	next func() float64
}

func NewCallbackOscillator(callback func() float64) *CallbackOscillator {
	return &CallbackOscillator{
		next: callback,
	}
}

func (o *CallbackOscillator) NextSample() float64 {
	return o.next()
}

func (o *CallbackOscillator) SetFreq(freq float64) {
	// No-op
}

func (o *CallbackOscillator) ResetPhase() {
	// No-op
}

func (o *CallbackOscillator) SetPhaseShift(phase float64) {
	// No-op
}
