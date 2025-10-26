package midi

import (
	"fmt"
	"math"

	"gitlab.com/gomidi/midi/v2"
)

type mode int

const (
	modeMenu mode = iota
	modePlay
)

const (
	wheelChannel = 1
	holdChannel  = 64
	forwardKey   = 50
	backwardKey  = 48
)

type Router struct {
	menu   *Menu
	voicer *Voicer

	mode mode
}

func NewRouter(voicer *Voicer, menu *Menu) *Router {
	return &Router{
		menu:   menu,
		voicer: voicer,
		mode:   modePlay,
	}
}

func (r *Router) Route(msg midi.Message) {
	var ch, key, vel uint8
	var channel, ctrl, val uint8

	switch {
	case msg.GetNoteStart(&ch, &key, &vel):
		if r.mode == modePlay {
			vel := math.Pow(float64(vel)/127.0, .8) // Exp curve TODO move to voicer
			r.voicer.NoteOn(KeysTable[key], vel)
			break
		}

		if key == forwardKey {
			r.menu.Forward()
			break
		}

		if key == backwardKey {
			r.menu.Backward()
			break
		}

	case msg.GetNoteEnd(&ch, &key):
		r.voicer.NoteOff(KeysTable[key]) // Always process note off to avoid stuck notes

	case msg.GetControlChange(&channel, &ctrl, &val):
		if channel == 0 && ctrl == holdChannel {
			if val >= 64 {
				r.mode = modeMenu
				break
			}
			r.mode = modePlay
		}

		if channel == 0 && ctrl == wheelChannel && r.mode == modeMenu {
			r.menu.Wheel(val)
		}
	default:
		fmt.Println("UNKNOWN COMMAND", msg.String())
	}
}
