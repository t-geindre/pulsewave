package audio

type Resettable interface {
	IsActive() bool
	Reset()
}
