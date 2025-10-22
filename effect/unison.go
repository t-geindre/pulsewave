package effect

import (
	"math"
	"math/rand"
	"synth/audio"
)

type Unison struct {
	*Merger
	makeVoice audio.SourceFactory
}

func NewUnison(makeVoice audio.SourceFactory, voices int, detuneCents, panSpread, curveExp float64) *Unison {
	u := &Unison{
		Merger:    NewMerger(),
		makeVoice: makeVoice,
	}

	// Pre-calc
	type w struct{ gL, gR float64 }
	weights := make([]w, voices)

	// t ∈ [-1..+1] sym around 0
	center := float64(voices-1) * 0.5

	for i := 0; i < voices; i++ {
		var t float64
		if voices == 1 {
			t = 0
		} else {
			t = (float64(i) - center) / center // [-1..+1]
		}

		// Exp curve
		curve := math.Copysign(math.Pow(math.Abs(t), curveExp), t)

		// Jitter (avoid phase locking)
		jitter := (rand.Float64()*2 - 1) * 0.05
		curveJ := curve * (1 + jitter)

		voice := NewTuner(makeVoice())

		if voices%2 == 1 && float64(i) == center {
			voice.SetDetuneCents(0)
			weights[i] = w{gL: math.Cos(math.Pi / 4), gR: math.Sin(math.Pi / 4)}
		} else {
			voice.SetDetuneCents(curveJ * detuneCents)

			// Equal-power pan with panSpread
			p := curve * panSpread           // [-panSpread..+panSpread]
			theta := (p + 1) * (math.Pi / 4) // -1 → 0, +1 → π/2
			weights[i] = w{gL: math.Cos(theta), gR: math.Sin(theta)}
		}

		u.Append(voice, 0, 0)
	}

	// ∑(gL²+gR²) = 1
	var esum float64
	for _, ww := range weights {
		esum += ww.gL*ww.gL + ww.gR*ww.gR
	}
	if esum == 0 {
		esum = 1
	}
	norm := 1.0 / math.Sqrt(esum)

	for i, ww := range weights {
		u.SetGain(i, ww.gL*norm, ww.gR*norm)
	}

	return u
}
