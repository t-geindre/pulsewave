package dsp

import "synth/audio"

type ConstParam struct {
	buf [audio.BlockSize]float32
}

func NewConstParam(value float32) *ConstParam {
	c := &ConstParam{}
	c.SetBase(value)
	return c
}

func (c *ConstParam) Resolve(uint64) []float32 {
	return c.buf[:] // No copy
}

func (c *ConstParam) SetBase(value float32) {
	for i := 0; i < audio.BlockSize; i++ {
		c.buf[i] = value
	}
}

func (c *ConstParam) ModInputs() *[]ParamModInput {
	panic("ConstParam param does not support modulations")
}
