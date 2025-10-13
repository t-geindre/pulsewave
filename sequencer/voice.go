package sequencer

type Voice interface {
	NoteOn(freq, velocity float64)
	NoteOff()
	NextSample() float64
	IsActive() bool
}

type VoiceFactory func() Voice
