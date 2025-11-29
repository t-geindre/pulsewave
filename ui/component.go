package ui

import (
	"synth/tree"

	"github.com/hajimehoshi/ebiten/v2"
)

type Component interface {
	Draw(*ebiten.Image)
	Update()
	Scroll(delta int)
	Focus()
	Blur()
	CurrentTarget() tree.Node // tree node
	// todo CurrentScroll() (node, max int) // node == -1 : no scroll bar
	// todo implement a focus function, especially for the preset node to be able to reset to load option
	// preventing accidental save
}
