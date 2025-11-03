package dsp

type Param interface {
	SetBase(value float32)
	GetBase() float32
	Resolve(cycle uint64) []float32
	ModInputs() *[]ParamModInput
}

type ParamSimple struct {
	base      float32
	inputs    []ParamModInput
	buf       [BlockSize]float32
	stampedAt uint64
}

func NewParam(base float32) *ParamSimple {
	return &ParamSimple{
		base: base,
	}
}

func (s *ParamSimple) SetBase(v float32)           { s.base = v }
func (s *ParamSimple) GetBase() float32            { return s.base }
func (s *ParamSimple) ModInputs() *[]ParamModInput { return &s.inputs }

func (s *ParamSimple) Resolve(cycle uint64) []float32 {
	if s.stampedAt == cycle {
		return s.buf[:]
	}

	for i := 0; i < BlockSize; i++ {
		s.buf[i] = s.base
	}

	for _, mi := range s.inputs {
		src := mi.Src.Resolve(cycle)
		if mi.Map == nil {
			for i := 0; i < BlockSize; i++ {
				s.buf[i] += mi.Amount * src[i]
			}
		} else {
			for i := 0; i < BlockSize; i++ {
				s.buf[i] += mi.Amount * mi.Map(src[i])
			}
		}
	}

	s.stampedAt = cycle

	return s.buf[:]
}

type ConstParam struct {
	buff [BlockSize]float32
}

func NewConstParam(v float32) *ConstParam {
	cp := &ConstParam{}
	for i := range cp.buff {
		cp.buff[i] = v
	}
	return cp
}

func (c *ConstParam) Resolve(cycle uint64) []float32 { return c.buff[:] }
func (c *ConstParam) SetBase(float32)                { panic("not implemented") } // const never changes
func (c *ConstParam) GetBase() float32               { panic("not implemented") } // const is not retrievable
func (c *ConstParam) ModInputs() *[]ParamModInput    { panic("not implemented") } // const never changes
