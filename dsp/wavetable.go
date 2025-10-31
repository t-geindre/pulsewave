package dsp

import (
	"math"
)

type Wavetable struct {
	Table []float32
	Size  int
}

func NewWavetable(size int, generator func(phase float64) float64) *Wavetable {
	table := make([]float32, size)
	for i := range table {
		phase := float64(i) / float64(size)
		table[i] = float32(generator(phase))
	}
	return &Wavetable{Table: table, Size: size}
}

func NewSineWavetable(size int) *Wavetable {
	return NewWavetable(size, func(phase float64) float64 {
		return math.Sin(2 * math.Pi * phase)
	})
}
