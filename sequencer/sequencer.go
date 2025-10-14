package sequencer

import (
	"math"
	"time"
)

type LoopMode int

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

func (s *Sequencer) NextSample() float64 {
	// No pattern, no sound
	if s.pattern == nil {
		return 0
	}

	// Avance step (manage possible multiple step boundaries)
	s.toNextStep--
	for s.toNextStep <= 0 {
		s.toNextStep += s.stepLength

		// Trigger all Notes scheduled at this step
		for {
			note := s.pattern.Next(s.step)
			if note == nil {
				break
			}
			s.triggerNote(note)
		}
		s.step++
	}

	sum := 0.0

	// Parcours à l'envers pour permettre la suppression en place
	for i := len(s.activeVoices) - 1; i >= 0; i-- {
		slot := s.activeVoices[i]

		// -- GATE / NoteOff (choisis le moment)
		// Si 'length' compte "échantillons restants y compris celui-ci",
		// tu peux décrémenter APRÈS la génération et noter off quand il devient 0.
		// Ici je le fais AVANT pour que le NoteOff tombe pile au bon sample suivant ta sémantique.
		if slot.length > 0 {
			slot.length--
			if slot.length == 0 {
				slot.voice.NoteOff()
			}
		}

		// -- Audio
		sum += slot.voice.NextSample()

		// -- Recycling: si l’enveloppe a fini, on libère la voix
		if !slot.voice.IsActive() {
			s.freeVoices = append(s.freeVoices, slot.voice)
			// remove i (stable mais O(n); acceptable vu peu de voix)
			s.activeVoices = append(s.activeVoices[:i], s.activeVoices[i+1:]...)
		}
	}

	return sum
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

func (s *Sequencer) IsActive() bool {
	return len(s.activeVoices) > 0 || (s.pattern != nil && s.pattern.IsActive())
}

func (s *Sequencer) Reset() {
	s.step = 0
	s.toNextStep = 0
	s.pattern.Reset()
	for _, voice := range s.activeVoices {
		voice.voice.NoteOff()
	}
	s.activeVoices = []*playingVoice{}
	s.freeVoices = []Voice{}
}
