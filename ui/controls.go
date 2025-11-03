package ui

import (
	"math"
	"synth/midi"
	"synth/msg"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Controls struct {
	midiIn             *msg.Queue
	upSince, downSince time.Time
}

func NewControls(midiIn *msg.Queue) *Controls {
	return &Controls{
		midiIn:    midiIn,
		upSince:   time.Time{},
		downSince: time.Time{},
	}
}

// Update poll inputs then returns forward, backward, scrollDelta
func (c *Controls) Update() (bool, bool, int) {
	fw, bw, scr := false, false, 0
	// todo implement repeat delay
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		fw = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		bw = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if c.upSince.IsZero() {
			c.upSince = time.Now()
			scr--
		} else {
			if time.Since(c.upSince) > 200*time.Millisecond {
				scr--
			}
			if time.Since(c.upSince) > 1000*time.Millisecond {
				scr--
			}
		}
	} else {
		c.upSince = time.Time{}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if c.downSince.IsZero() {
			c.downSince = time.Now()
			scr++
		} else {
			if time.Since(c.downSince) > 200*time.Millisecond {
				scr++
			}
			if time.Since(c.downSince) > 1000*time.Millisecond {
				scr++
			}
		}
	} else {
		c.downSince = time.Time{}
	}
	c.midiIn.Drain(10, func(m msg.Message) {
		switch m.Kind {
		case midi.ControlChangeKind: // Todo controller config
			if m.Key == 112 && m.Val8 != 0 {
				v := 0.0
				if m.Val8 < 64 {
					v += float64(m.Val8 - 64)
				} else {
					v -= float64(64 - m.Val8)
				}
				sh := math.Pow(v+1, 2) / 4 // make it less sensitive
				sh = math.Copysign(sh, v)
				scr += int(v)
			}
		default:
			// ignore
		}
	})

	return fw, bw, scr
}
