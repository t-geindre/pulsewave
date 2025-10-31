package ui

import (
	"math"
	"synth/midi"
	"synth/msg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Controls struct {
	midiIn *msg.Queue
}

func NewControls(midiIn *msg.Queue) *Controls {
	return &Controls{
		midiIn: midiIn,
	}
}

// Update poll inputs then returns forward, backward, scrollDelta
func (c *Controls) Update() (bool, bool, int) {
	fw, bw, scr := false, false, 0
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		fw = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		bw = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		scr++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		scr--
	}
	c.midiIn.Drain(10, func(m msg.Message) {
		switch m.Kind {
		case midi.ControlChangeKind:
			if m.Key == 112 && m.Val8 != 0 {
				v := 0.0
				if m.Val8 > 64 {
					v += float64(m.Val8 - 64)
				} else {
					v -= float64(64 - m.Val8)
				}
				sh := math.Pow(v+1, 2) / 3 // make it less sensitive
				sh = math.Copysign(sh, v)
				scr += int(v)
			}
		default:
			// ignore
		}
	})

	return fw, bw, scr
}
