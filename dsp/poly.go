package dsp

import "synth/audio"

type PolyVoice struct {
	voices map[int]*Voice
	inputs map[int]*Input
	*Mixer
}

func NewPolyVoice(voices int, factory func() *Voice) *PolyVoice {
	p := &PolyVoice{
		voices: make(map[int]*Voice, voices),
		inputs: make(map[int]*Input, voices),
		Mixer:  NewMixer(NewParam(1), false),
	}

	gp, pp := NewParam(1), NewParam(0)
	for i := 0; i < voices; i++ {
		vc := factory()
		in := NewInput(vc, gp, pp)
		in.Mute = true
		p.Mixer.Add(in)
		p.voices[i], p.inputs[i] = vc, in
	}

	return p
}

func (p *PolyVoice) NoteOn(key int, vel float32) {
	if vc, ok := p.voices[key]; ok {
		vc.NoteOn(key, vel)
		return
	}

	for n, vc := range p.voices {
		if vc.IsIdle() {
			in := p.inputs[n]

			delete(p.voices, n)
			delete(p.inputs, n)

			p.voices[key], p.inputs[key] = vc, in

			vc.NoteOn(key, vel)
			return
		}
	}
}

func (p *PolyVoice) NoteOff(key int) {
	if vc, ok := p.voices[key]; ok {
		vc.NoteOff()
	}
}

func (p *PolyVoice) Process(b *audio.Block) {
	for i, vc := range p.voices {
		p.inputs[i].Mute = vc.IsIdle()
	}

	p.Mixer.Process(b)
}
