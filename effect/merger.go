package effect

import (
	"math"
	"synth/audio"
)

type Merger struct {
	srcs  []audio.Source
	gains [][2]float64 // [left, right]
}

func NewMerger() *Merger {
	return &Merger{}
}

func (m *Merger) NextValue() (float64, float64) {
	vl, vr := .0, .0

	for i, src := range m.srcs {
		l, r := src.NextValue()
		vl += l * m.gains[i][0]
		vr += r * m.gains[i][1]
	}

	return vl, vr
}

func (m *Merger) Append(src audio.Source, left, right float64) {
	m.srcs = append(m.srcs, src)
	m.gains = append(m.gains, [2]float64{left, right})

	// Normalize gains if their sum of absolute values exceeds 1
	for i := range m.gains {
		s := .0
		for j := range m.gains[i] {
			s += math.Abs(m.gains[i][j])
		}
		if s > 1 {
			for j := range m.gains[i] {
				m.gains[i][j] /= s
			}
		}

	}
}

func (m *Merger) IsActive() bool {
	for _, src := range m.srcs {
		if src.IsActive() {
			return true
		}
	}
	return false
}

func (m *Merger) Reset() {
	for _, src := range m.srcs {
		src.Reset()
	}
}

func (m *Merger) NoteOn(freq, velocity float64) {
	for _, src := range m.srcs {
		src.NoteOn(freq, velocity)
	}
}

func (m *Merger) NoteOff() {
	for _, src := range m.srcs {
		src.NoteOff()
	}
}

func (m *Merger) SetGain(index int, left, right float64) {
	if index < 0 || index >= len(m.gains) {
		return
	}

	m.gains[index][0] = left
	m.gains[index][1] = right
}
