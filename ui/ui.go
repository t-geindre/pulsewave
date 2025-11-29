package ui

import (
	"synth/assets"
	"synth/msg"
	"synth/tree"

	"github.com/hajimehoshi/ebiten/v2"
)

// Todo get it fron config
const (
	BodyWidth  = 375
	BodyHeight = 211
	BodyStartX = 52
	BodyStartY = 70
	SlideSpeed = .08
)

type Ui struct {
	background *ebiten.Image
	w, h       int
	controls   Controls

	messenger *msg.Messenger

	components Components
	current    tree.Node
	next       tree.Node

	nextTrans float64
	transDir  float64

	transLeft    *ebiten.Image
	transRight   *ebiten.Image
	bodyClipMask *ebiten.Image

	ftps *Ftps
}

func NewUi(
	asts *assets.Loader,
	messenger *msg.Messenger,
	ctrls Controls,
	cmps Components,
	root tree.Node,
) (*Ui, error) {
	// BG + window size accordingly
	bg, err := asts.GetImage("ui/background")
	if err != nil {
		return nil, err
	}

	bds := bg.Bounds()
	ebiten.SetWindowSize(bds.Dx(), bds.Dy())

	return &Ui{
		background:   bg,
		w:            bds.Dx(),
		h:            bds.Dy(),
		components:   cmps,
		messenger:    messenger,
		current:      root,
		controls:     ctrls,
		bodyClipMask: ebiten.NewImage(BodyWidth, BodyHeight),
		transLeft:    ebiten.NewImage(BodyWidth, BodyHeight),
		transRight:   ebiten.NewImage(BodyWidth, BodyHeight),
	}, nil
}

func (u *Ui) Update() error {
	u.messenger.Process()

	// Transitioning
	if u.next != nil {
		u.components[u.next].Update()
		u.nextTrans += SlideSpeed
		if u.nextTrans > 1 {
			u.current = u.next
			u.next = nil
			u.nextTrans = 0
		}
		return nil
	}

	hDelta, vDelta := u.controls.Update()

	// Forward
	if hDelta > 0 {
		tr := u.components[u.current].CurrentTarget()

		// Target node, follow redirect if any
		if target, ok := tr.(tree.RedirectionNode); ok {
			redirect := target.GetRedirection()
			if redirect != nil {
				tr = redirect
			}
		}

		if tr != nil {
			u.next = tr
			u.components[u.current].Blur()
			u.components[u.next].Focus()
			u.transDir = 1
			if tr == u.current.Parent() {
				// Component sending us back up the tree should slide left
				u.transDir = -1
			}
		}
		return nil

	}

	// Backward
	if hDelta < 0 {
		// Leaving a subtree
		leaving := u.current.Context().IsLeavingSubTree(u.current)
		if leaving != nil {
			u.next = leaving
			u.transDir = -1
			return nil
		}

		// Normal back to parent
		p := u.current.Parent()
		if p != nil {
			u.next = p
			u.transDir = -1
		}
		return nil
	}

	// Vertical scroll
	if vDelta != 0 {
		u.components[u.current].Scroll(vDelta)
	}
	u.components[u.current].Update()

	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {
	screen.DrawImage(u.background, nil)
	u.bodyClipMask.Clear()

	if u.next != nil {
		// Transitioning
		ease := easeOutCubic(u.nextTrans)
		if u.transDir == -1 {
			ease = 1 - ease
		}

		u.transLeft.Clear()
		u.components[u.current].Draw(u.transLeft)

		u.transRight.Clear()
		u.components[u.next].Draw(u.transRight)

		lOpts := ebiten.DrawImageOptions{}
		lOpts.GeoM.Translate(-ease*BodyWidth, 0)

		rOpts := ebiten.DrawImageOptions{}
		rOpts.GeoM.Translate(BodyWidth-ease*BodyWidth, 0)

		if u.transDir == 1 {
			u.bodyClipMask.DrawImage(u.transLeft, &lOpts)
			u.bodyClipMask.DrawImage(u.transRight, &rOpts)
		} else {
			u.bodyClipMask.DrawImage(u.transLeft, &rOpts)
			u.bodyClipMask.DrawImage(u.transRight, &lOpts)
		}
	} else {
		// Normal draw
		u.components[u.current].Draw(u.bodyClipMask)
	}

	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(BodyStartX, BodyStartY)
	screen.DrawImage(u.bodyClipMask, ops)

	if u.ftps != nil {
		u.ftps.Draw(screen)
	}
}

func (u *Ui) Layout(_, _ int) (int, int) {
	return u.w, u.h
}

func (u *Ui) ToggleFtpsDisplay() {
	if u.ftps == nil {
		u.ftps = NewFtps()
		return
	}
	u.ftps = nil
}
