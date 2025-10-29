package dsp

import "math"

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

// panGains computes left and right gain factors for a given pan
// value p in [-1, +1] using equal-power panning
func panGains(p float32) (gl, gr float32) {
	theta := (float64(p) + 1.0) * 0.25 * math.Pi
	gl = float32(math.Cos(theta))
	gr = float32(math.Sin(theta))
	return
}

// softClip applies a soft clipping function to the input x
func softClip(x float32) float32 {
	ax := float32(math.Abs(float64(x)))
	return x / (1 + ax)
}

// LPF (Uni pole) from cutoff (Hz)
func lpfCoef(cutHz, sr float64) float32 {
	if cutHz <= 0 {
		return 0
	}

	// bilinear simple / one-pole exp
	// alpha ~= 1 - exp(-2π fc / sr)
	return float32(1 - math.Exp(-2*math.Pi*cutHz/sr))
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
