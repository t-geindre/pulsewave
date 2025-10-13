package sequencer

import (
	"math"
	"time"
)

type LoopMode int

const (
	LoopOff    LoopMode = iota
	LoopStrict          // repeat right after last noteOff
	LoopSoft            // wait for all notes to finish before repeating
)

type playingVoice struct {
	voice  Voice
	length int
}

type Sequencer struct {
	// Voices
	freeVoices   []Voice
	activeVoices []*playingVoice
	maxVoices    int
	voiceFactory VoiceFactory

	// Pattern
	pattern *Pattern
	loop    LoopMode

	// Timing
	stepLength   int
	toNextStep   int
	step         int
	sr           float64
	stepsPerBeat int
}

func NewSequencer(sampleRate float64, tempo float64, maxVoices int, stepsPerBeat int, voiceFactory VoiceFactory) *Sequencer {
	return &Sequencer{
		maxVoices:    maxVoices,
		voiceFactory: voiceFactory,
		freeVoices:   []Voice{},
		loop:         LoopOff,
		stepLength: int(math.Round(
			60.0 * float64(sampleRate) / (tempo * float64(stepsPerBeat)),
		)),
		toNextStep:   0,
		step:         0,
		sr:           sampleRate,
		stepsPerBeat: stepsPerBeat,
	}
}

func (s *Sequencer) SetPattern(p *Pattern) {
	s.pattern = p
}

func (s *Sequencer) SetLoopMode(mode LoopMode) {
	s.loop = mode
}

func (s *Sequencer) NextSample() float64 {
	// No pattern, no sound
	if s.pattern == nil {
		return 0
	}

	// Advance step
	s.toNextStep -= 1
	if s.toNextStep <= 0 {
		// Handle looping
		if !s.pattern.IsActive() {
			switch s.loop {
			case LoopStrict:
				s.pattern.Reset()
				s.step = 0
			case LoopSoft:
				if len(s.activeVoices) == 0 {
					s.pattern.Reset()
					s.step = 0
				}
			default:
			}
		}

		s.toNextStep += s.stepLength
		// Trigger all notes at this step
		for {
			note := s.pattern.Next(s.step)
			if note == nil {
				break
			}
			s.triggerNote(note)
		}
		s.step++
	}

	// Mix active voices
	v := 0.0
	for i, voice := range s.activeVoices {
		v += voice.voice.NextSample()

		// Remove finished voices
		if !voice.voice.IsActive() {
			s.freeVoices = append(s.freeVoices, voice.voice)
			s.activeVoices = append(s.activeVoices[:i], s.activeVoices[i+1:]...)
		}
		// Off notes
		if voice.length <= 0 {
			continue
		}

		voice.length--
		if voice.length <= 0 {
			voice.voice.NoteOff()
		}
	}

	return v
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

func (s *Sequencer) getFreeVoice() Voice {
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
