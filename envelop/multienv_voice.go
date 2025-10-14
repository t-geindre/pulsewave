package envelop

type MultiEnvVoice struct {
	inner    *Voice
	envelops []Envelop
}

func NewMultiEnvVoice(inner *Voice, envelops ...Envelop) *MultiEnvVoice {
	return &MultiEnvVoice{
		inner:    inner,
		envelops: envelops,
	}
}

func (m *MultiEnvVoice) NoteOn(freq, velocity float64) {
	for _, env := range m.envelops {
		env.NoteOn()
	}
	m.inner.NoteOn(freq, velocity)
}

func (m *MultiEnvVoice) NoteOff() {
	for _, env := range m.envelops {
		env.NoteOff()
	}
	m.inner.NoteOff()
}

func (m *MultiEnvVoice) NextSample() float64 {
	return m.inner.NextSample()
}

func (m *MultiEnvVoice) IsActive() bool {
	for _, env := range m.envelops {
		if env.IsActive() {
			return true
		}
	}

	return m.inner.IsActive()
}
