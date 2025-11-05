package ui

import (
	"math"
	"synth/midi"
	"synth/msg"
	"time"
)

type MidiControls struct {
	midiIn   *msg.Queue
	lastHor  time.Time
	dblClick time.Duration
}

func NewMidiControls(midiIn *msg.Queue) *MidiControls {
	return &MidiControls{
		midiIn:   midiIn,
		dblClick: 200 * time.Millisecond,
	}
}

func (c *MidiControls) Update() (int, int) {
	horDelta, vertDelta := 0, 0

	c.midiIn.Drain(20, func(m msg.Message) {
		switch m.Kind {
		case midi.ControlChangeKind: // Todo controller config
			if m.Key == 112 && m.Val8 != 0 {
				v := 0
				if m.Val8 < 64 {
					v += int(64 - m.Val8)
				} else {
					v -= int(m.Val8 - 64)
				}
				vertDelta += v
			}
			if m.Key == 113 && m.Val8 != 127 { // 0 = release, 127 = press
				if c.lastHor.IsZero() {
					c.lastHor = time.Now()
					return
				}
				if !c.lastHor.IsZero() && time.Since(c.lastHor) < c.dblClick {
					c.lastHor = time.Time{}
					horDelta = -1
				}
			}
		default:
			// ignore
		}
	})

	if !c.lastHor.IsZero() && time.Since(c.lastHor) > c.dblClick {
		c.lastHor = time.Time{}
		horDelta = 1
	}

	dir := 1
	if vertDelta < 0 {
		dir = -1
	}

	vertDelta = int(math.Pow(float64(vertDelta), 4)) * dir

	return horDelta, vertDelta
}
