package dsp

type Param interface {
	SetBase(value float32)
	Resolve(cycle uint64) []float32
	ModInputs() *[]ParamModInput
}
