package preset

import "synth/msg"

const AudioSource msg.Source = 10

const ParamUpdateKind msg.Kind = 1
const ParamPullAllKind msg.Kind = 2

const (
	// Unison parameters
	UnisonOnOff = iota // 0 = off, 1 = on
	UnisonPanSpread
	UnisonPhaseSpread
	UnisonDetuneSpread
	UnisonCurveGamma
	UnisonVoices

	// Pitch Modulation
	PitchLfoOnOff
	PitchLfoAmount
	PitchLfoShape
	PitchLfoFreq
	PitchLfoPhase
	PitchAdsrOnOff
	PitchAdsrAmount
	PitchAdsrAttack
	PitchAdsrDecay
	PitchAdsrSustain
	PitchAdsrRelease

	// Feedback Delay parameters
	FBOnOff // 0 = off, 1 = on
	FBDelayParam
	FBFeedBack
	FBMix
	FBTone

	// Amp Envelope parameters
	AmpEnvAttack
	AmpEnvDecay
	AmpEnvSustain
	AmpEnvRelease

	// Oscillator parameters
	Osc0Shape
	Osc0Detune
	Osc0Gain
	Osc0Phase
	Osc0Pw

	Osc1Shape
	Osc1Detune
	Osc1Gain
	Osc1Phase
	Osc1Pw

	Osc2Shape
	Osc2Detune
	Osc2Gain
	Osc2Phase
	Osc2Pw

	// Low Pass Filter parameters
	LPFOnOff
	LPFCutoff
	LPFResonance

	LpfLfoOnOff
	LpfLfoAmount
	LpfLfoShape
	LpfLfoFreq
	LpfLfoPhase

	LpfAdsrOnOff
	LpfAdsrAmount
	LpfAdsrAttack
	LpfAdsrDecay
	LpfAdsrSustain
	LpfAdsrRelease
)

/*
Message
	Source Source AudioSource
	Kind   Kind   ParamUpdateKind
	key    uint8  0-255 parameter ID
	ValF  float32
}
*/
