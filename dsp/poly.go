package dsp

type polyVoice struct {
	key   int
	voice *Voice
	input *Input
	index uint64
}

const (
	PolyStealOldest = iota
	PolyStealLowest
	PolyStealHighest
)

type PolyVoice struct {
	*Mixer
	voices       []*polyVoice
	index        uint64
	stealMode    Param
	activeVoices Param
}

func NewPolyVoice(maxVoices int, activeVoices Param, stealMode Param, factory func() *Voice) *PolyVoice {
	p := &PolyVoice{
		voices:       make([]*polyVoice, maxVoices),
		Mixer:        NewMixer(nil, false),
		stealMode:    stealMode,
		activeVoices: activeVoices,
	}

	for i := 0; i < maxVoices; i++ {
		vc := &polyVoice{}
		vc.voice = factory()
		vc.input = NewInput(vc.voice, nil, nil) // no pan/gain, mixer fast path
		vc.input.Mute = true
		p.Mixer.Add(vc.input)
		p.voices[i] = vc
	}

	return p
}

func (p *PolyVoice) NoteOn(key int, vel float32) {
	p.index++
	av := int(p.activeVoices.Resolve(p.index)[0])

	// Same key retrigger
	for i := 0; i < av; i++ {
		v := p.voices[i]
		if v.key == key {
			v.index = p.index
			v.voice.NoteOn(key, vel)
			v.input.Mute = false
			return
		}
	}

	// Find idle voice
	for i := 0; i < av; i++ {
		v := p.voices[i]
		if v.voice.IsIdle() {
			v.key = key
			v.index = p.index
			v.voice.NoteOn(key, vel)
			v.input.Mute = false
			return
		}
	}

	// Steal voice
	var lru *polyVoice
	switch int(p.stealMode.Resolve(p.index)[0]) {
	case PolyStealLowest:
		lru = p.stealLowest(av)
	case PolyStealHighest:
		lru = p.stealHighest(av)
	default:
		lru = p.stealOldest(av)
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
			s.input.Mute = true
			s.key = 0
			s.index = 0
			break
		}
	}
}

func (p *PolyVoice) Process(b *Block) {
	for _, s := range p.voices {
		s.input.Mute = s.voice.IsIdle()
	}
	p.Mixer.Process(b)
}

func (p *PolyVoice) AllNotesOff() {
	for _, s := range p.voices {
		s.voice.NoteOff()
		s.input.Mute = true
		s.key = 0
		s.index = 0
	}
}

func (p *PolyVoice) stealOldest(av int) *polyVoice {
	var lru *polyVoice
	for i := 0; i < av; i++ {
		s := p.voices[i]
		if lru == nil || s.index < lru.index {
			lru = s
		}
	}
	return lru
}

func (p *PolyVoice) stealLowest(av int) *polyVoice {
	var lru *polyVoice
	for i := 0; i < av; i++ {
		s := p.voices[i]
		if lru == nil || s.key < lru.key {
			lru = s
		}
	}
	return lru
}

func (p *PolyVoice) stealHighest(av int) *polyVoice {
	var lru *polyVoice
	for i := 0; i < av; i++ {
		s := p.voices[i]
		if lru == nil || s.key > lru.key {
			lru = s
		}
	}
	return lru
}
