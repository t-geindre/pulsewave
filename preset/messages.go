package preset

import "synth/msg"

// LoadSavePresetKind msg.key = preset slot, msg.val8 = 0 load, 1 save
const LoadSavePresetKind msg.Kind = 21

// UpdateParameterKind msg.key = parameter ID, msg.valF = parameter value
const UpdateParameterKind msg.Kind = 20

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

	// Voices
	VoicesStealMode  = 55
	VoicesActive     = 56
	VoicesPitchGlide = 57

	// Noise oscillator
	NoiseGain = 58
	NoiseType = 62

	// Sub oscillator
	SubOscShape     = 59
	SubOscGain      = 60
	SubOscTranspose = 61

	// LFOs
	Lfo0rate  = 63
	Lfo0Phase = 64
	Lfo0Shape = 65

	Lfo1rate  = 66
	Lfo1Phase = 67
	Lfo1Shape = 68

	Lfo2rate  = 69
	Lfo2Phase = 70
	Lfo2Shape = 71

	// ADSRs
	Adsr0Attack  = 72
	Adsr0Decay   = 73
	Adsr0Sustain = 74
	Adsr0Release = 75

	Adsr1Attack  = 76
	Adsr1Decay   = 77
	Adsr1Sustain = 78
	Adsr1Release = 79

	Adsr2Attack  = 80
	Adsr2Decay   = 81
	Adsr2Sustain = 82
	Adsr2Release = 83

	// No parameter
	None = 255
)

// ModulationUpdateKind msg.key = slot, msg.channel = source, msg.val8 = destination, msg.valF = amount, msg.val16 = shape
const ModulationUpdateKind msg.Kind = 22

const (
	ModSlots       = 10 // number of modulation slots
	ModKeysSpacing = 10 // x/ModKeysSpacing = slot, x%ModKeysSpacing = param

	ModParamSrc = 0
	ModParamDst = 1
	ModParamAmt = 2
	ModParamShp = 3
)

const (
	ModSrcVelocity = 0
	ModSrcLfo1     = 1
	ModSrcLfo2     = 2
	ModSrcLfo3     = 3
	ModSrcAdsr1    = 4
	ModSrcAdsr2    = 5
	ModSrcAdsr3    = 6
)

const (
	ModShapeLinear      = 0
	ModShapeExponential = 1
	ModShapeLogarithmic = 2
)
