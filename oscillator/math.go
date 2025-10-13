package oscillator

import "math"

func frac(x float64) float64 { return x - math.Floor(x) }

func frac01(x float64) float64 {
	x = math.Mod(x, 1.0)
	if x < 0 {
		x += 1.0
	}
	return x
}

func polyBLEP(t, dt float64) float64 {
	if t < dt {
		t /= dt
		return t + t - t*t - 1.0 // up
	}
	if t > 1.0-dt {
		t = (t - 1.0) / dt
		return t*t + 2.0*t + 1.0 // down
	}
	return 0.0
}
