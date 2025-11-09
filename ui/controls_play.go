package ui

import "github.com/hajimehoshi/ebiten/v2"

type playKey struct {
	key  ebiten.Key
	note uint8
	down bool
}

type PlayControls struct {
	keys      []*playKey
	messenger *Messenger
}

func NewPlayControls(messenger *Messenger) *PlayControls {
	return &PlayControls{
		keys: []*playKey{
			{key: ebiten.KeyA, note: 60},         // C4
			{key: ebiten.KeyS, note: 62},         // D4
			{key: ebiten.KeyD, note: 64},         // E4
			{key: ebiten.KeyF, note: 65},         // F4
			{key: ebiten.KeyG, note: 67},         // G4
			{key: ebiten.KeyH, note: 69},         // A4
			{key: ebiten.KeyJ, note: 71},         // B4
			{key: ebiten.KeyK, note: 72},         // C5
			{key: ebiten.KeyL, note: 74},         // D5
			{key: ebiten.KeySemicolon, note: 76}, // E5
			{key: ebiten.KeyW, note: 61},         // C#4
			{key: ebiten.KeyE, note: 63},         // D#4
			{key: ebiten.KeyT, note: 66},         // F#4
			{key: ebiten.KeyY, note: 68},         // G#4
			{key: ebiten.KeyU, note: 70},         // A#4
			{key: ebiten.KeyO, note: 73},         // C#5
			{key: ebiten.KeyP, note: 75},         // D#5
		},
		messenger: messenger,
	}
}

func (p *PlayControls) Update() (horDelta, vertDelta int) {
	for _, pk := range p.keys {
		if pk.down {
			if !ebiten.IsKeyPressed(pk.key) {
				p.messenger.NoteOff(0, pk.note)
				pk.down = false
			}
			continue
		}
		if ebiten.IsKeyPressed(pk.key) {
			p.messenger.NoteOn(0, pk.note, 100)
			pk.down = true
		}
	}

	return 0, 0
}
