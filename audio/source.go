package audio

type Source interface {
	NextSample() float64
}
