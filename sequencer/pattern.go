package sequencer

import "sort"

type NoteSpec struct {
	At       int
	Freq     float64
	Length   int
	Velocity float64
	Gate     float64
}

type Pattern struct {
	dirty bool
	notes []NoteSpec
	index int
	next  int
}

func NewPattern() *Pattern {
	return &Pattern{
		dirty: true,
		notes: []NoteSpec{},
		index: 0,
	}
}

func (p *Pattern) AddAt(at int, freq float64, length int, velocity, gate float64) {
	if at+length > p.next {
		p.next = at + length
	}

	if freq <= 0 || length <= 0 {
		return
	}

	p.notes = append(p.notes, NoteSpec{At: at, Freq: freq, Length: length, Velocity: velocity, Gate: gate})
	p.dirty = true
}

func (p *Pattern) Append(freq float64, length int, velocity, gate float64) {
	p.AddAt(p.next, freq, length, velocity, gate)
}

func (p *Pattern) Next(at int) *NoteSpec {
	if len(p.notes) == 0 || p.index >= len(p.notes) || p.notes[p.index].At > at {
		return nil
	}
	p.sort()
	note := &p.notes[p.index]
	p.index++

	return note
}

func (p *Pattern) IsActive() bool {
	return p.index < len(p.notes)
}

func (p *Pattern) Reset() {
	p.index = 0
}

func (p *Pattern) sort() {
	if p.dirty {
		sort.Slice(p.notes, func(i, j int) bool {
			return p.notes[i].At < p.notes[j].At
		})
		p.dirty = false
	}
}
