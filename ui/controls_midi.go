package ui

import (
	"math"
	"time"
)

type MidiControls struct {
	lastHor             time.Time
	dblClick            time.Duration
	vertDelta, horDelta int
}

func NewMidiControls() *MidiControls {
	return &MidiControls{
		dblClick: 200 * time.Millisecond,
	}
}

func (c *MidiControls) Update() (int, int) {
	horDelta, vertDelta := c.horDelta, c.vertDelta
	c.horDelta, c.vertDelta = 0, 0
	return horDelta, vertDelta
}

func (c *MidiControls) ControlChange(key, val uint8) {
	if key == 112 && val != 0 {
		v := 0
		if val < 64 {
			v += int(64 - val)
		} else {
			v -= int(val - 64)
		}
		c.vertDelta += v
	}
	if key == 113 && val != 127 { // 0 = release, 127 = press
		if c.lastHor.IsZero() {
			c.lastHor = time.Now()
			return
		}
		if !c.lastHor.IsZero() && time.Since(c.lastHor) < c.dblClick {
			c.lastHor = time.Time{}
			c.horDelta = -1
		}
	}

	if !c.lastHor.IsZero() && time.Since(c.lastHor) > c.dblClick {
		c.lastHor = time.Time{}
		c.horDelta = 1
	}

	dir := 1
	if c.vertDelta < 0 {
		dir = -1
	}

	c.vertDelta = int(math.Pow(float64(c.vertDelta), 4)) * dir
}
