package ui

func easeOutCubic(t float64) float64 {
	return t * (3 - 3*t + t*t)
}
