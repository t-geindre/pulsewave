package dsp

import (
	"github.com/viterin/vek/vek32"
)

type Input struct {
	Src  Node
	Gain Param // linear
	Pan  Param // -1..+1 (0 center), equal-power
	Mute bool
}

// NewInput creates a new Input node.
// src: source Node to mix.
// gain: per-sample gain (can be nil and SHOULD BE if not required for better performance).
// pan: per-sample pan (-1..+1, can be nil and SHOULD BE if not required for better performance).
func NewInput(src Node, gain, pan Param) *Input {
	return &Input{Src: src, Gain: gain, Pan: pan}
}

type Mixer struct {
	Inputs []*Input

	MasterGain Param // linear
	SoftClip   bool  // soft limiter/tanh

	accL [BlockSize]float32
	accR [BlockSize]float32

	tmp Block
}

// NewMixer creates a new Mixer.
// masterGain: overall gain applied after mixing all inputs (can be nil).
// softClip: if true, applies soft clipping to the output.
func NewMixer(masterGain Param, softClip bool) *Mixer {
	m := &Mixer{MasterGain: masterGain, SoftClip: softClip}
	m.tmp.L = [BlockSize]float32{}
	m.tmp.R = [BlockSize]float32{}
	return m
}

func (m *Mixer) Add(in *Input) int {
	m.Inputs = append(m.Inputs, in)
	return len(m.Inputs) - 1
}

func (m *Mixer) Clear() { m.Inputs = m.Inputs[:0] }

func (m *Mixer) Process(b *Block) {
	for i := range m.accL {
		m.accL[i] = 0
		m.accR[i] = 0
	}

	m.tmp.Cycle = b.Cycle

	for _, in := range m.Inputs {
		if in == nil || in.Src == nil || in.Mute {
			continue
		}

		var gainB, panB []float32
		if in.Gain != nil {
			gainB = in.Gain.Resolve(b.Cycle)
		}
		if in.Pan != nil {
			panB = in.Pan.Resolve(b.Cycle)
		}

		in.Src.Process(&m.tmp)

		// Fast path: no gain, no pan
		if gainB == nil && panB == nil {
			var tmpL, tmpR [BlockSize]float32
			vek32.Add_Into(tmpL[:], m.accL[:], m.tmp.L[:])
			vek32.Add_Into(tmpR[:], m.accR[:], m.tmp.R[:])
			copy(m.accL[:], tmpL[:])
			copy(m.accR[:], tmpR[:])
			continue
		}

		// Slow path: per-sample gain/pan
		var gL, gR [BlockSize]float32
		for i := 0; i < BlockSize; i++ {
			g := float32(1)
			if gainB != nil {
				g = gainB[i]
			}
			gl, gr := float32(1), float32(1)
			if panB != nil {
				gl, gr = fastPanGains(panB[i])
			}
			gL[i], gR[i] = g*gl, g*gr
		}

		var tmpL, tmpR, mixL, mixR [BlockSize]float32
		vek32.Mul_Into(tmpL[:], m.tmp.L[:], gL[:])
		vek32.Mul_Into(tmpR[:], m.tmp.R[:], gR[:])
		vek32.Add_Into(mixL[:], m.accL[:], tmpL[:])
		vek32.Add_Into(mixR[:], m.accR[:], tmpR[:])
		copy(m.accL[:], mixL[:])
		copy(m.accR[:], mixR[:])
	}

	// Apply master gain / soft clip
	if m.MasterGain != nil {
		gb := m.MasterGain.Resolve(b.Cycle)
		if m.SoftClip {
			var tmpL, tmpR [BlockSize]float32
			vek32.Mul_Into(tmpL[:], m.accL[:], gb[:])
			vek32.Mul_Into(tmpR[:], m.accR[:], gb[:])
			for i := range b.L {
				b.L[i] = softClip(tmpL[i])
				b.R[i] = softClip(tmpR[i])
			}
		} else {
			vek32.Mul_Into(b.L[:], m.accL[:], gb[:])
			vek32.Mul_Into(b.R[:], m.accR[:], gb[:])
		}
		return
	}

	if m.SoftClip {
		for i := range b.L {
			b.L[i] = softClip(m.accL[i])
			b.R[i] = softClip(m.accR[i])
		}
	} else {
		copy(b.L[:], m.accL[:])
		copy(b.R[:], m.accR[:])
	}
}

func (m *Mixer) Reset() {
	for _, in := range m.Inputs {
		if in != nil && in.Src != nil {
			in.Src.Reset()
		}
	}
}
