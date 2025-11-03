package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Component interface {
	Draw(*ebiten.Image)
	Update()
	Scroll(delta int)
	CurrentTarget() Node // tree node
	// CurrentScroll() (node, max int) // node == -1 : no scroll bar
}
