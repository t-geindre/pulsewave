package dsp

type ParamModulator interface {
	Resolve(cycle uint64) []float32 // read-only, len == BlockSize
}

// ParamModInput = une modulation vers un Param
type ParamModInput struct {
	Src    ParamModulator          // source
	Amount float32                 // gain (+/-)
	Map    func(x float32) float32 // courbe optionnelle (ex: expo), appliquée AVANT Amount
}

func NewModInput(src ParamModulator, amount float32, mapping func(x float32) float32) ParamModInput {
	return ParamModInput{
		Src:    src,
		Amount: amount,
		Map:    mapping,
	}
}
