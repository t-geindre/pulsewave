package dsp

type ParamModulator interface {
	Resolve(cycle uint64) []float32 // read-only, len == BlockSize
}

type ParamModInput interface {
	Src() ParamModulator
	Amount() Param
	Map() func(x float32) float32
}

type ParamModulatorInput struct {
	src    ParamModulator
	amount Param
	mapf   func(x float32) float32
}

func NewModInput(src ParamModulator, amount Param, mapping func(x float32) float32) ParamModInput {
	return ParamModulatorInput{
		src:    src,
		amount: amount,
		mapf:   mapping,
	}
}

func (p ParamModulatorInput) Src() ParamModulator {
	return p.src
}

func (p ParamModulatorInput) Amount() Param {
	return p.amount
}
func (p ParamModulatorInput) Map() func(x float32) float32 {
	return p.mapf
}
