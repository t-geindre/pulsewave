package ui

import (
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
)

type Component interface {
	Draw(*ebiten.Image)
	Update()
	Scroll(delta int)
	CurrentTarget() *preset.Node // tree node
	// CurrentScroll() (node, max int) // node == -1 : no scroll bar
}
