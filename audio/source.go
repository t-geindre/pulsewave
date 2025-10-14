package audio

type Source interface {
	NextSample() float64
	IsActive() bool
	Reset()
}
