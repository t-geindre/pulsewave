package ui

import (
	"synth/midi"
	"synth/msg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Controls struct {
	midiIn *msg.Queue

	fw, bw     bool
	more, less int
}

func NewControls(midiIn *msg.Queue) *Controls {
	return &Controls{
		midiIn: midiIn,
	}
}

func (c *Controls) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		c.fw = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		c.bw = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		c.more++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		c.less++
	}
	c.midiIn.Drain(10, func(m msg.Message) {
		switch m.Kind {
		case midi.ControlChangeKind:
			if m.Key == 112 && m.Val != 0 {
				if m.Val > 64 {
					c.more += int(m.Val - 64)
				} else {
					c.less += int(64 - m.Val)
				}
			}
		default:
			// ignore
		}
	})
}

func (c *Controls) Consume() (forward, back bool, more, less int) {
	forward = c.fw
	back = c.bw
	more = c.more
	less = c.less

	c.fw = false
	c.bw = false
	c.more = 0
	c.less = 0

	return
}
