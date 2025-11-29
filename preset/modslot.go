package preset

import "synth/dsp"

type ModSlot struct {
	Source           uint8
	Destination      uint8
	Amount           float32
	Shape            uint8
	GlobalModInput   dsp.ParamModInput
	PerVoiceModInput []dsp.ParamModInput
}
