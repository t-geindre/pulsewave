package dsp

import "synth/audio"

type Node interface {
	audio.Source
	Reset()
}
