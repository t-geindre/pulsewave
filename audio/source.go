package audio

type Source interface {
	// Process fills the given Block with audio data
	Process(*Block)
}
