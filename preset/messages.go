package preset

import "synth/msg"

const PresetUpdateKind msg.Kind = 20
const LoadSavePresetKind msg.Kind = 21

const (
	// Unison parameters
	UnisonOnOff        = 0 // 0 = off, 1 = on
	UnisonPanSpread    = 1
	UnisonPhaseSpread  = 2
	UnisonDetuneSpread = 3
	UnisonCurveGamma   = 4
	UnisonVoices       = 5

	// Pitch Modulation
	PitchLfoOnOff    = 6
	PitchLfoAmount   = 7
	PitchLfoShape    = 8
	PitchLfoFreq     = 9
	PitchLfoPhase    = 10
	PitchAdsrOnOff   = 11
	PitchAdsrAmount  = 12
	PitchAdsrAttack  = 13
	PitchAdsrDecay   = 14
	PitchAdsrSustain = 15
	PitchAdsrRelease = 16

	// Feedback Delay parameters
	FBOnOff      = 17 // 0 = off, 1 = on
	FBDelayParam = 18
	FBFeedBack   = 19
	FBMix        = 20
	FBTone       = 21

	// Amp Envelope parameters
	AmpEnvAttack  = 22
	AmpEnvDecay   = 23
	AmpEnvSustain = 24
	AmpEnvRelease = 25

	// Oscillator parameters
	Osc0Shape  = 26
	Osc0Detune = 27
	Osc0Gain   = 28
	Osc0Phase  = 29
	Osc0Pw     = 30

	Osc1Shape  = 31
	Osc1Detune = 32
	Osc1Gain   = 33
	Osc1Phase  = 34
	Osc1Pw     = 35

	Osc2Shape  = 36
	Osc2Detune = 37
	Osc2Gain   = 38
	Osc2Phase  = 39
	Osc2Pw     = 40

	// Low Pass Filter parameters
	LPFOnOff     = 41
	LPFCutoff    = 42
	LPFResonance = 43

	LpfLfoOnOff  = 44
	LpfLfoAmount = 45
	LpfLfoShape  = 46
	LpfLfoFreq   = 47
	LpfLfoPhase  = 48

	LpfAdsrOnOff   = 49
	LpfAdsrAmount  = 50
	LpfAdsrAttack  = 51
	LpfAdsrDecay   = 52
	LpfAdsrSustain = 53
	LpfAdsrRelease = 54

	// None, assigned to nothing
	NONE = 255
)

/*
Message
	Source Source AudioSource
	Kind   Kind   PresetUpdateKind
	key    uint8  0-255 parameter ID
	ValF  float32
}
*/
