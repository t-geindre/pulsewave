package envelop

type Envelop interface {
	NoteOn()
	NoteOff()
	Next() float64
	Value() float64
	IsActive() bool
}
