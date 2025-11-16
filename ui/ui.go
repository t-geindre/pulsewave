package ui

import (
	"fmt"
	"synth/assets"
	"synth/msg"
	"synth/tree"

	"github.com/hajimehoshi/ebiten/v2"
)

var ErrorUnknownNodeType = fmt.Errorf("unknown node type")

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

	components map[tree.Node]Component
	current    tree.Node
	next       tree.Node
	nextTrans  float64
	transDir   float64
	tree       tree.Node

	transLeft    *ebiten.Image
	transRight   *ebiten.Image
	bodyClipMask *ebiten.Image

	ftps *Ftps
}

func NewUi(asts *assets.Loader, messenger *msg.Messenger, audiQ *AudioQueue, presets []string) (*Ui, error) {
	// BG + window size accordingly
	bg, err := asts.GetImage("ui/background")
	if err != nil {
		return nil, err
	}
	bds := bg.Bounds()

	ebiten.SetWindowSize(bds.Dx(), bds.Dy())
	ebiten.SetWindowTitle("Pulsewave")

	// Menu tree + controls
	menu := tree.NewTree(presets)
	menu.Attach(messenger)

	ctrls := NewMultiControls(
		NewPlayControls(messenger),
		NewKeyboardControls(),
		NewTouchControls(),
	)

	ui := &Ui{
		background:   bg,
		w:            bds.Dx(),
		h:            bds.Dy(),
		components:   make(map[tree.Node]Component),
		messenger:    messenger,
		current:      menu,
		tree:         menu,
		controls:     ctrls,
		bodyClipMask: ebiten.NewImage(BodyWidth, BodyHeight),
		transLeft:    ebiten.NewImage(BodyWidth, BodyHeight),
		transRight:   ebiten.NewImage(BodyWidth, BodyHeight),
	}

	err = ui.buildComponents(asts, ui.current, audiQ)
	if err != nil {
		return nil, err
	}

	return ui, nil
}

func (u *Ui) Update() error {
	u.messenger.Process()

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

	if hDelta > 0 {
		tr := u.components[u.current].CurrentTarget()
		if tr != nil {
			u.next = tr
			u.transDir = 1
			if tr == u.current.Parent() {
				// Component sending us back up the tree should slide left
				u.transDir = -1
			}
		}
		return nil

	}
	if hDelta < 0 && u.current.Parent() != nil {
		u.next = u.current.Parent()
		u.transDir = -1
		return nil
	}

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

func (u *Ui) buildComponents(asts *assets.Loader, n tree.Node, aq *AudioQueue) error {
	switch node := n.(type) {
	case tree.SliderNode:
		comp, err := NewSlider(asts, node)
		if err != nil {
			return err
		}
		u.components[node] = comp
	case tree.SelectorNode:
		comp, err := NewSelector(asts, node)
		if err != nil {
			return err
		}
		u.components[node] = comp
	case tree.FeatureNode:
		switch node.Feature() {
		case tree.FeatureOscilloscope:
			comp, err := NewOscilloscope(aq, 16384)
			if err != nil {
				return err
			}
			u.components[node] = comp
			// todo add default case that errors out
		}
	default:
		if len(node.Children()) > 0 {
			comp, err := NewList(asts, node)
			if err != nil {
				return err
			}
			u.components[node] = comp
			break
		}
		// return ErrorUnknownNodeType todo decide if we want to error out here
	}

	for _, child := range n.Children() {
		err := u.buildComponents(asts, child, aq)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Ui) ToggleFpsDisplay() {
	if u.ftps == nil {
		u.ftps = NewFtps()
		return
	}
	u.ftps = nil
}
