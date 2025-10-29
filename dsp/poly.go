package dsp

type PolyVoice struct {
	voices  map[int]*Voice
	max     int
	factory func() *Voice
	*Mixer
}

func NewPolyVoice(max int, factory func() *Voice) *PolyVoice {
	return &PolyVoice{
		voices:  make(map[int]*Voice),
		max:     max,
		factory: factory,
		Mixer:   NewMixer(NewParam(1), true),
	}
}

func (p *PolyVoice) NoteOn(key int, vel float32) {
	if vc, ok := p.voices[key]; ok {
		vc.NoteOn(key, vel)
		return
	}

	if len(p.voices) >= p.max {
		for n, vc := range p.voices {
			if vc.IsIdle() {
				delete(p.voices, n)
				p.voices[key] = vc
				vc.NoteOn(key, vel)
				return
			}
		}
		return // Skip
	}

	vc := p.factory()
	vc.NoteOn(key, vel)
	p.voices[key] = vc

	p.Mixer.Add(NewInput(vc, NewParam(1), NewParam(0)))
}

func (p *PolyVoice) NoteOff(key int) {
	if vc, ok := p.voices[key]; ok {
		vc.NoteOff()
	}
}
