package audio

type CallbackSrc struct {
	call func() (float64, float64)
}

func NewCallbackSrc(call func() (float64, float64)) *CallbackSrc {
	return &CallbackSrc{
		call: call,
	}
}

func (c *CallbackSrc) NextValue() (float64, float64) {
	return c.call()
}

func (c *CallbackSrc) IsActive() bool {
	return false
}

func (c *CallbackSrc) Reset() {
}

func (c *CallbackSrc) NoteOn(freq, velocity float64) {
}

func (c *CallbackSrc) NoteOff() {
}
