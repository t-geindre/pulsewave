package preset

import "synth/dsp"

type Manager struct {
	current int
	presets []*Polysynth
}

func (p *Manager) NoteOn(i int, f float32) {
}

func (p *Manager) NoteOff(i int) {
}

func (p *Manager) SetPitchBend(st float32) {
}

func (p *Manager) Process(block *dsp.Block) {
}

func (p *Manager) Reset(soft bool) {
}
