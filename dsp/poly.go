package dsp

type polyVoice struct {
	key   int
	voice *Voice
	input *Input
	index uint64
}

type PolyVoice struct {
	*Mixer
	voices []*polyVoice
	index  uint64
}

func NewPolyVoice(voices int, factory func() *Voice) *PolyVoice {
	p := &PolyVoice{
		voices: make([]*polyVoice, voices),
		Mixer:  NewMixer(NewParam(1), false),
	}

	gp, pp := NewParam(1), NewParam(0)
	for i := 0; i < voices; i++ {
		vc := &polyVoice{}
		vc.voice = factory()
		vc.input = NewInput(vc.voice, gp, pp)
		vc.input.Mute = true
		p.Mixer.Add(vc.input)
		p.voices[i] = vc
	}

	return p
}

func (p *PolyVoice) NoteOn(key int, vel float32) {
	p.index++

	for _, s := range p.voices {
		if s.key == key {
			s.index = p.index
			s.voice.NoteOn(key, vel)
			s.input.Mute = false
			return
		}
	}

	for _, s := range p.voices {
		if s.voice.IsIdle() {
			s.key = key
			s.index = p.index
			s.voice.NoteOn(key, vel)
			s.input.Mute = false
			return
		}
	}

	lru := p.voices[0]
	for _, s := range p.voices[1:] {
		if s.index < lru.index {
			lru = s
		}
	}

	lru.key = key
	lru.index = p.index
	lru.voice.NoteOff()
	lru.voice.NoteOn(key, vel)
	lru.input.Mute = false
}

func (p *PolyVoice) NoteOff(key int) {
	for _, s := range p.voices {
		if s.key == key {
			s.voice.NoteOff()
			return
		}
	}
}

func (p *PolyVoice) Process(b *Block) {
	for _, s := range p.voices {
		s.input.Mute = s.voice.IsIdle()
	}
	p.Mixer.Process(b)
}
