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
	dirty  bool
	Notes  []NoteSpec
	index  int
	next   int
	length int
}

func NewPattern() *Pattern {
	return &Pattern{
		dirty: true,
		Notes: []NoteSpec{},
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

	p.Notes = append(p.Notes, NoteSpec{At: at, Freq: freq, Length: length, Velocity: velocity, Gate: gate})
	p.dirty = true
}

func (p *Pattern) Append(freq float64, length int, velocity, gate float64) {
	p.AddAt(p.next, freq, length, velocity, gate)
}

func (p *Pattern) Prepend(freq float64, length int, velocity, gate float64) {
	p.AddAt(0, freq, length, velocity, gate)

	for i := range p.Notes {
		p.Notes[i].At += length
	}
	p.next += length
	p.dirty = true
}

func (p *Pattern) Next(at int) *NoteSpec {
	if len(p.Notes) == 0 || p.index >= len(p.Notes) || p.Notes[p.index].At > at {
		return nil
	}
	p.sort()
	note := &p.Notes[p.index]
	p.index++

	return note
}

func (p *Pattern) IsActive() bool {
	return p.index < len(p.Notes)
}

func (p *Pattern) Reset() {
	p.index = 0
}

func (p *Pattern) Clone() *Pattern {
	clone := NewPattern()

	clone.Notes = make([]NoteSpec, len(p.Notes))
	copy(clone.Notes, p.Notes)

	clone.next = p.next
	clone.dirty = p.dirty
	clone.length = p.length
	clone.index = p.index

	return clone
}

func (p *Pattern) Length() int {
	p.sort()
	return p.length
}

func (p *Pattern) Move(offset int) {
	for i := range p.Notes {
		p.Notes[i].At += offset
	}
	p.next += offset
	p.dirty = true
}

func (p *Pattern) sort() {
	if p.dirty {
		sort.Slice(p.Notes, func(i, j int) bool {
			return p.Notes[i].At < p.Notes[j].At
		})
		p.dirty = false

		p.length = 0
		for _, n := range p.Notes {
			p.length += n.Length
		}
	}
}
