package midi

import (
	"synth/dsp"
	"synth/msg"
)

type Instrument interface {
	NoteOn(int, float32)
	NoteOff(int)
}

type Player struct {
	dsp.Node
	inst  Instrument
	queue *msg.Queue
}

func NewPlayer(src dsp.Node, inst Instrument, queue *msg.Queue) *Player {
	return &Player{
		Node:  src,
		inst:  inst,
		queue: queue,
	}
}

func (p *Player) Process(block *dsp.Block) {
	var m msg.Message
	for p.queue.TryRead(&m) {
		switch m.Type {
		case NoteOnKind:
			// Todo handle vel properly with LUT (precalculated curve)
			p.inst.NoteOn(int(m.V1), float32(m.V2)/127)
		case NoteOffKind:
			p.inst.NoteOff(int(m.V1))
		}
	}

	p.Node.Process(block)
}

func (p *Player) Reset() {
	p.Node.Reset()
}
