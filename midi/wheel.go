package midi

import (
	"math"
	"time"
)

type Wheel struct {
	lastT time.Time
	lastV uint8

	// réglages
	Deadband int     // ex: 2 ou 3
	Kquad    float64 // poids de l'accélération (0 = linéaire pur). ex: 0.15
	Wrap127  bool    // true si 0↔127 “wrap”

	// état
	carry float64 // accumulateur fractionnaire
}

func NewWheel() *Wheel {
	return &Wheel{
		Deadband: 1,
		Kquad:    0.15,
		Wrap127:  false,
	}
}

func (w *Wheel) Update(v uint8) int {
	now := time.Now()
	if w.lastT.IsZero() || time.Since(w.lastT) > time.Millisecond*50 {
		w.lastT, w.lastV = now, v
		return 0
	}

	// delta brut
	raw := int(v) - int(w.lastV)
	// gestion optionnelle du wrap 0↔127
	if w.Wrap127 {
		if raw > 64 {
			raw -= 128
		} else if raw < -64 {
			raw += 128
		}
	}

	// deadband : on ignore les très petits mouvements
	absRaw := abs(raw)
	if absRaw <= w.Deadband {
		w.lastT, w.lastV = now, v
		return 0
	}

	// amplitude “utile” (au-delà du deadband)
	mag := float64(absRaw - w.Deadband)
	sign := 1.0
	if raw < 0 {
		sign = -1.0
	}

	// courbe simple : linéaire + petit terme quadratique
	// => lent: ~1 par unité ; rapide: monte progressivement
	stepsCont := (mag * (1 + w.Kquad*mag)) * sign

	// on accumule pour ne pas perdre les fractions
	w.carry += stepsCont
	out := int(math.Round(w.carry))
	w.carry -= float64(out)

	w.lastT, w.lastV = now, v
	return out
}
