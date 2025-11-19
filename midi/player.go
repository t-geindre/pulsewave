package midi

import (
	"synth/msg"
	"synth/settings"
)

type Instrument interface {
	NoteOn(int, float32)
	NoteOff(int)
	SetPitchBend(st float32)
}

type Player struct {
	inst        Instrument
	pitchBendSt float32
}

func NewPlayer(inst Instrument) *Player {
	return &Player{
		inst: inst,
	}
}

func (p *Player) HandleMessage(m msg.Message) {
	switch m.Kind {
	case NoteOnKind:
		// Todo handle vel properly with LUT (precalculated curve)
		p.inst.NoteOn(int(m.Key), float32(m.Val8)/127)
	case NoteOffKind:
		p.inst.NoteOff(int(m.Key))
	case PitchBendKind:
		rel := float32(0)
		if m.Val16 >= 128 || m.Val16 <= -128 {
			rel = float32(m.Val16) / 8192.0 * p.pitchBendSt
		}
		p.inst.SetPitchBend(rel)
	case settings.SettingUpdateKind:
		if m.Key == settings.PitchBendRange {
			p.pitchBendSt = m.ValF
		}
	}
}
