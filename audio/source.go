package audio

type Source interface {
	// NextValue returns the next left and right audio samples.
	NextValue() (L, R float64)

	// IsActive returns whether the source is still active (producing sound).
	IsActive() bool

	// Reset resets the source to its initial state.
	Reset()

	// NoteOn starts playing a note with the given frequency and velocity.
	NoteOn(freq, velocity float64)

	// NoteOff stops playing the current note.
	NoteOff()
}

type SourceFactory func() Source
