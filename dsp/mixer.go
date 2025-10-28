package dsp

import (
	"synth/audio"
)

type Input struct {
	Src  audio.Source
	Gain Param // linear
	Pan  Param // -1..+1 (0 center), equal-power
	Mute bool
}

type Mixer struct {
	Inputs []*Input

	MasterGain Param // linear
	SoftClip   bool  // soft limiter/tanh

	accL [audio.BlockSize]float32
	accR [audio.BlockSize]float32

	tmp audio.Block
}

func NewMixer(masterGain Param, softClip bool) *Mixer {
	m := &Mixer{
		MasterGain: masterGain,
		SoftClip:   softClip,
	}

	m.tmp.L = [audio.BlockSize]float32{}
	m.tmp.R = [audio.BlockSize]float32{}

	return m
}

func (m *Mixer) Add(in *Input) { m.Inputs = append(m.Inputs, in) }

func (m *Mixer) Process(b *audio.Block) {
	for i := 0; i < audio.BlockSize; i++ {
		// Todo check if zeroing is needed
		m.accL[i] = 0
		m.accR[i] = 0
	}

	// Sum
	m.tmp.Cycle = b.Cycle
	for _, in := range m.Inputs {
		if in == nil || in.Src == nil {
			continue
		}
		if in.Mute {
			continue
		}

		// Pull param buffers
		var gainB, panB []float32
		if in.Gain != nil {
			gainB = in.Gain.Resolve(b.Cycle)
		}
		if in.Pan != nil {
			panB = in.Pan.Resolve(b.Cycle)
		}

		// Pull block
		in.Src.Process(&m.tmp)

		// Sum
		for i := 0; i < audio.BlockSize; i++ {
			g := float32(1)
			if gainB != nil {
				g = gainB[i]
			}
			gl, gr := float32(1), float32(1)
			if panB != nil {
				gl, gr = panGains(panB[i])
			}
			m.accL[i] += m.tmp.L[i] * g * gl
			m.accR[i] += m.tmp.R[i] * g * gr
		}
	}

	// Master gain
	if m.MasterGain != nil {
		gb := m.MasterGain.Resolve(b.Cycle)
		for i := 0; i < audio.BlockSize; i++ {
			yL := m.accL[i] * gb[i]
			yR := m.accR[i] * gb[i]
			if m.SoftClip {
				yL = softClip(yL)
				yR = softClip(yR)
			}
			b.L[i] = yL
			b.R[i] = yR
		}
		return
	}

	if m.SoftClip {
		for i := 0; i < audio.BlockSize; i++ {
			b.L[i] = softClip(m.accL[i])
			b.R[i] = softClip(m.accR[i])
		}
	} else {
		copy(b.L[:], m.accL[:])
		copy(b.R[:], m.accR[:])
	}
}
