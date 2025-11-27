package dsp

type Param interface {
	SetBase(value float32)
	GetBase() float32
	Resolve(cycle uint64) []float32
	ModInputs() *[]ParamModInput
	RemoveModInputBySource(ParamModulator) bool // returns true if something was removed
	AddModInput(ParamModInput)
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
		src := mi.Src().Resolve(cycle)
		amount := mi.Amount().Resolve(cycle)
		mapf := mi.Map()
		if mapf == nil {
			for i := 0; i < BlockSize; i++ {
				s.buf[i] += amount[i] * src[i]
			}
		} else {
			for i := 0; i < BlockSize; i++ {
				s.buf[i] += amount[i] * mapf(src[i])
			}
		}
	}

	s.stampedAt = cycle

	return s.buf[:]
}

func (s *ParamSimple) RemoveModInputBySource(src ParamModulator) bool {
	removed := false
	newInputs := s.inputs[:0]
	for _, mi := range s.inputs {
		if mi.Src() != src {
			newInputs = append(newInputs, mi)
		} else {
			removed = true
		}
	}
	s.inputs = newInputs

	return removed
}

func (s *ParamSimple) AddModInput(mi ParamModInput) {
	s.inputs = append(s.inputs, mi)
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
func (c *ConstParam) RemoveModInputBySource(ParamModulator) bool {
	panic("not implemented")
}
func (c *ConstParam) AddModInput(ParamModInput) {
	panic("not implemented")
}
