package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type playKey struct {
	key  ebiten.Key
	note uint8
	down bool
}

type PlayControls struct {
	keys      []*playKey
	messenger *Messenger
	oct       uint8
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
			{key: ebiten.KeyZ, note: 36},         // C2
			{key: ebiten.KeyX, note: 38},         // D2
			{key: ebiten.KeyC, note: 40},         // E2
			{key: ebiten.KeyV, note: 41},         // F2
			{key: ebiten.KeyB, note: 43},         // G2
			{key: ebiten.KeyN, note: 45},         // A2
			{key: ebiten.KeyM, note: 47},         // B2
			{key: ebiten.KeyComma, note: 48},     // C3
			{key: ebiten.KeyPeriod, note: 50},    // D3
			{key: ebiten.KeySlash, note: 52},     // E3
		},
		messenger: messenger,
	}
}

func (p *PlayControls) Update() (horDelta, vertDelta int) {
	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		p.oct++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		p.oct--
	}
	for _, pk := range p.keys {
		if pk.down {
			if !ebiten.IsKeyPressed(pk.key) {
				p.messenger.NoteOff(0, pk.note+p.oct*12)
				pk.down = false
			}
			continue
		}
		if ebiten.IsKeyPressed(pk.key) {
			p.messenger.NoteOn(0, pk.note+p.oct*12, 100)
			pk.down = true
		}
	}

	return 0, 0
}
