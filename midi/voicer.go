package midi

import (
	"sync"
	"synth/audio"
)

type Voicer struct {
	max     int
	voices  map[float64]audio.Source
	factory audio.SourceFactory
	mut     sync.Mutex
}

func NewVoicer(maxVoices int, factory audio.SourceFactory) *Voicer {
	v := &Voicer{
		max:     maxVoices,
		factory: factory,
	}
	v.Reset()
	return v
}

func (v *Voicer) NextValue() (L, R float64) {
	v.mut.Lock()
	defer v.mut.Unlock()

	for _, src := range v.voices {
		l, r := src.NextValue()
		L += l
		R += r
	}

	return L, R
}

func (v *Voicer) IsActive() bool {
	v.mut.Lock()
	defer v.mut.Unlock()

	for _, src := range v.voices {
		if src.IsActive() {
			return true
		}
	}

	return false
}

func (v *Voicer) Reset() {
	v.voices = make(map[float64]audio.Source, v.max)
}

func (v *Voicer) NoteOn(freq, velocity float64) {
	v.mut.Lock()
	defer v.mut.Unlock()

	acc, ok := v.voices[freq]
	if ok {
		acc.NoteOn(freq, velocity)
		return
	}

	if len(v.voices) >= v.max {
		for k, src := range v.voices {
			if !src.IsActive() {
				delete(v.voices, k)
				src.NoteOn(freq, velocity)
				v.voices[freq] = src
				return
			}
		}
		// all voices are active, skip
		return
	}

	src := v.factory()
	src.NoteOn(freq, velocity)

	v.voices[freq] = src
}

func (v *Voicer) NoteOff(freq float64) {
	v.mut.Lock()
	defer v.mut.Unlock()

	if acc, ok := v.voices[freq]; ok {
		acc.NoteOff(freq)
	}
}
