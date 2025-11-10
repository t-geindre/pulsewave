package midi

import (
	"synth/dsp"
	"synth/msg"
)

type Instrument interface {
	NoteOn(int, float32)
	NoteOff(int)
	SetPitchBend(st float32)
}

type Player struct {
	dsp.Node
	inst  Instrument
	queue *msg.Queue
	msg   msg.Message
}

func NewPlayer(src dsp.Node, inst Instrument, queue *msg.Queue) *Player {
	return &Player{
		Node:  src,
		inst:  inst,
		queue: queue,
	}
}

func (p *Player) Process(block *dsp.Block) {
	for p.queue.TryRead(&p.msg) {
		p.processMessage(p.msg)
	}

	p.Node.Process(block)
}

func (p *Player) Reset(soft bool) {
	p.Node.Reset(soft)
}

func (p *Player) processMessage(m msg.Message) {
	switch m.Kind {
	case NoteOnKind:
		// Todo handle vel properly with LUT (precalculated curve)
		p.inst.NoteOn(int(m.Key), float32(m.Val8)/127)
	case NoteOffKind:
		p.inst.NoteOff(int(m.Key))
	case PitchBendKind:
		rel := float32(0)
		if m.Val16 >= 128 || m.Val16 <= -128 {
			rel = float32(m.Val16) / 8192.0 * 4.0 // 4 semitones range
		}
		p.inst.SetPitchBend(rel)
	}
}
