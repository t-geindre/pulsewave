package sequencer

import (
	"math"
	"synth/audio"
	"time"
)

type playingVoice struct {
	voice  audio.Source
	length int
}

type pattern struct {
	at      int
	pattern Pattern
}

type Sequencer struct {
	// Voices
	freeVoices   []audio.Source
	activeVoices []*playingVoice
	maxVoices    int
	voiceFactory audio.SourceFactory

	// Pattern
	patterns []*Pattern
	index    int

	// Timing
	stepLength   int
	toNextStep   int
	step         int
	sr           float64
	stepsPerBeat int
}

func NewSequencer(sampleRate float64, tempo float64, maxVoices int, stepsPerBeat int, voiceFactory audio.SourceFactory) *Sequencer {
	return &Sequencer{
		maxVoices:    maxVoices,
		voiceFactory: voiceFactory,
		freeVoices:   []audio.Source{},
		stepLength: int(math.Round(
			60.0 * float64(sampleRate) / (tempo * float64(stepsPerBeat)),
		)),
		toNextStep:   0,
		step:         0,
		sr:           sampleRate,
		stepsPerBeat: stepsPerBeat,
		patterns:     make([]*Pattern, 0),
	}
}

func (s *Sequencer) Append(p *Pattern) {
	p = p.Clone()
	p.Move(s.index)
	s.index += p.Length()
	s.patterns = append(s.patterns, p)
}

func (s *Sequencer) AppendAndRepeat(p *Pattern, times int) {
	for i := 0; i < times; i++ {
		s.Append(p)
	}
}

func (s *Sequencer) NextValue() (float64, float64) {
	if len(s.patterns) == 0 {
		return 0, 0
	}

	// Advance step
	s.toNextStep--
	for s.toNextStep <= 0 {
		s.toNextStep += s.stepLength

		// Trigger all Notes scheduled at this step
		for {
			triggered := false
			for _, p := range s.patterns {
				note := p.Next(s.step)
				if note != nil {
					s.triggerNote(note)
					triggered = true
				}
			}
			if !triggered {
				break
			}
		}
		s.step++
	}

	vl, vr := .0, .0

	// Reverse iteration to allow removal
	for i := len(s.activeVoices) - 1; i >= 0; i-- {
		slot := s.activeVoices[i]

		if slot.length > 0 {
			slot.length--
			if slot.length == 0 {
				slot.voice.NoteOff()
			}
		}

		l, r := slot.voice.NextValue()
		vl += l
		vr += r

		if !slot.voice.IsActive() {
			s.freeVoices = append(s.freeVoices, slot.voice)
			s.activeVoices = append(s.activeVoices[:i], s.activeVoices[i+1:]...)
		}
	}

	return vl, vr
}

func (s *Sequencer) triggerNote(note *NoteSpec) {
	// Get a voice
	voice := s.getFreeVoice()
	if voice == nil {
		// No free voice available
		return
	}

	// Start the voice
	voice.NoteOn(note.Freq, note.Velocity)

	// AddAt to active voices
	s.activeVoices = append(s.activeVoices, &playingVoice{
		voice:  voice,
		length: int(math.Round(float64(note.Length*s.stepLength) * note.Gate)),
	})
}

func (s *Sequencer) getFreeVoice() audio.Source {
	if len(s.freeVoices) > 0 {
		voice := s.freeVoices[len(s.freeVoices)-1]
		s.freeVoices = s.freeVoices[:len(s.freeVoices)-1]
		return voice
	}

	if len(s.activeVoices) < s.maxVoices {
		return s.voiceFactory()
	}

	return nil
}

func (s *Sequencer) GetStepDuration() time.Duration {
	return time.Duration(s.stepLength*int(time.Second)) / time.Duration(s.sr)
}

func (s *Sequencer) GetBeatDuration() time.Duration {
	return time.Duration(s.stepLength*s.stepsPerBeat*int(time.Second)) / time.Duration(s.sr)
}

func (s *Sequencer) IsActive() bool {
	if len(s.patterns) == 0 {
		return false
	}

	if len(s.activeVoices) > 0 {
		return true
	}

	for _, p := range s.patterns {
		if p.IsActive() {
			return true
		}
	}

	return false
}

func (s *Sequencer) Reset() {
	s.step = 0
	s.toNextStep = 0
	for _, p := range s.patterns {
		p.Reset()
	}
	for _, voice := range s.activeVoices {
		voice.voice.NoteOff()
	}
	s.activeVoices = []*playingVoice{}
	s.freeVoices = []audio.Source{}
}

func (s *Sequencer) NoteOn(_, _ float64) {
}

func (s *Sequencer) NoteOff() {
}
