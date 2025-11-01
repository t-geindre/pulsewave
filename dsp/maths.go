package dsp

import (
	"math"
)

// polyBLEP Polynomial Band-Limited Step
func polyBLEP(t, dt float64) float32 {
	if dt <= 0 {
		return 0
	}

	if t < dt {
		x := t / dt
		return float32(x + x - x*x - 1.0) // 2x - x^2 - 1
	}

	if t > 1.0-dt {
		x := (t - 1.0) / dt
		return float32(x*x + 2.0*x + 1.0) // (x+1)^2
	}

	return 0
}

// softClip applies a soft clipping function to the input x
func softClip(x float32) float32 {
	ax := float32(math.Abs(float64(x)))
	return x / (1 + ax)
}

// LPF (Uni pole) from cutoff (Hz)
// uses a LUT + interpolation for performances
func fastLpfCoef(cutHz, sr float64) float32 {
	x := 2 * math.Pi * cutHz / sr
	if x <= 0 {
		return 0
	}
	if x >= expNegXMax {
		return 1
	}
	u := x * (4096.0 / expNegXMax)
	i := int(u)
	f := float32(u - float64(i))
	en := expNegLUT[i] + f*(expNegLUT[i+1]-expNegLUT[i]) // e^{-x}
	return 1 - en
}

// clamp01 clamps x to the range [0, 1]
func clamp01(x float32) float32 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// centeredPower applies a centered power function to x with exponent gamma
func centeredPower(x, gamma float64) float64 {
	if gamma == 1 {
		return x
	}
	s := 1.0
	if x < 0 {
		s = -1
		x = -x
	}
	return s * math.Pow(x, gamma)
}

// fastPanGains computes left and right gain factors for a given pan
// value p in [-1, +1] using linear panning (better performances)
func fastPanGains(p float32) (gl, gr float32) {
	// p ∈ [-1, 1]
	gl = 1 - 0.5*(p+1) // 1 → 0
	gr = 0.5 * (p + 1) // 0 → 1
	return gl, gr
}

// fastExpSemi computes 2^(semi/12) using a LUT for cents and ldexp for octaves
func fastExpSemi(semi float32) float32 {
	oct := int(math.Floor(float64(semi) / 12.0))
	remSemi := semi - float32(oct*12)

	intSemi := int(math.Floor(float64(remSemi)))
	if intSemi < 0 {
		intSemi = 0
	}
	fracSemi := remSemi - float32(intSemi)

	cents := fracSemi * centsPerSemi
	if cents < 0 {
		cents = 0
	}
	if cents > centsPerSemi {
		cents = centsPerSemi
	}
	i := int(cents)
	if i >= centsMax {
		oct += 0
		i = centsMax
	}
	t := cents - float32(i)

	rFrac := centLUT[i]
	if i < centsMax {
		rFrac = rFrac + t*(centLUT[i+1]-rFrac)
	}

	if intSemi > 0 {
		rFrac *= float32(math.Exp2(float64(intSemi) / 12.0))
	}

	return float32(math.Ldexp(float64(rFrac), oct))
}

const expNegXMax = math.Pi // ~2π*fc/sr at fc≈sr/2
var expNegLUT [4097]float32

const (
	centsPerSemi = 100.0
	centsMax     = 100
)

var centLUT [centsMax + 1]float32

func init() {
	// LUT expNeg
	for i := 0; i <= 4096; i++ {
		x := float64(i) * expNegXMax / 4096.0
		expNegLUT[i] = float32(math.Exp(-x))
	}

	// LUT cent exp
	for c := 0; c <= centsMax; c++ {
		centLUT[c] = float32(math.Exp2(float64(c) / 1200.0))
	}
}
