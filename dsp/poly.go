package dsp

type polyVoice struct {
	key   int
	voice *Voice
	input *Input
	index uint64
	gate  bool
}

const (
	PolyStealOldest = iota
	PolyStealLowest
	PolyStealHighest
)

const MaxStolenRetain = 16

type PolyVoice struct {
	*Mixer
	voices       []*polyVoice
	index        uint64
	stealMode    Param
	activeVoices Param

	stolen     [MaxStolenRetain]int
	stolenHead int
	stolenSize int
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
			v.gate = true
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
			v.gate = true
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

	if lru.gate {
		p.enqueueStolen(lru.key)
	}

	lru.key = key
	lru.index = p.index
	lru.voice.NoteOff()
	lru.voice.NoteOn(key, vel)
	lru.input.Mute = false
	lru.gate = true
}

func (p *PolyVoice) NoteOff(key int) {
	if p.dropStolen(key) {
		return
	}

	for _, s := range p.voices {
		if s.key == key {
			s.voice.NoteOff()
			s.input.Mute = true
			s.key = 0
			s.index = 0
			s.gate = false

			// Re-trigger stolen note if any
			if stolenKey, found := p.dequeueStolen(); found {
				s.key = stolenKey
				s.index = p.index
				s.voice.NoteOn(stolenKey, 1.0) // velocity hardcoded to 1.0
				s.input.Mute = false
				s.gate = true
			}
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

func (p *PolyVoice) enqueueStolen(key int) {
	if p.stolenSize >= MaxStolenRetain {
		return
	}
	pos := (p.stolenHead + p.stolenSize) % MaxStolenRetain
	p.stolen[pos] = key
	p.stolenSize++
}

func (p *PolyVoice) dequeueStolen() (int, bool) {
	if p.stolenSize == 0 {
		return 0, false
	}

	top := (p.stolenHead + p.stolenSize - 1 + MaxStolenRetain) % MaxStolenRetain
	key := p.stolen[top]
	p.stolenSize--

	return key, true
}

func (p *PolyVoice) dropStolen(key int) bool {
	for i := 0; i < p.stolenSize; i++ {
		idx := (p.stolenHead + i) % MaxStolenRetain
		if p.stolen[idx] == key {
			for j := i; j < p.stolenSize-1; j++ {
				from := (p.stolenHead + j + 1) % MaxStolenRetain
				to := (p.stolenHead + j) % MaxStolenRetain
				p.stolen[to] = p.stolen[from]
			}
			p.stolenSize--
			return true
		}
	}
	return false
}
