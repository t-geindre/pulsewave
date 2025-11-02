package ui

import (
	"synth/assets"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
)

// Todo get it fron config
const (
	BodyHeight = 211
	BodyWidth  = 376
	BodyStartX = 51
	BodyStartY = 70
	SlideSpeed = .08
)

type Ui struct {
	background *ebiten.Image
	w, h       int
	controls   *Controls

	components map[*preset.Node]Component
	current    *preset.Node
	next       *preset.Node
	nextTrans  float64
	transDir   float64

	transLeft    *ebiten.Image
	transRight   *ebiten.Image
	bodyClipMask *ebiten.Image
}

func NewUi(asts *assets.Loader, ctrl *Controls, root *preset.Node) (*Ui, error) {
	// BG + window size accordingly
	bg, err := asts.GetImage("ui/background")
	if err != nil {
		return nil, err
	}
	bds := bg.Bounds()
	ebiten.SetWindowSize(bds.Dx(), bds.Dy())

	ui := &Ui{
		background:   bg,
		w:            bds.Dx(),
		h:            bds.Dy(),
		components:   make(map[*preset.Node]Component),
		current:      root,
		controls:     ctrl,
		bodyClipMask: ebiten.NewImage(BodyWidth, BodyHeight),
		transLeft:    ebiten.NewImage(BodyWidth, BodyHeight),
		transRight:   ebiten.NewImage(BodyWidth, BodyHeight),
	}

	err = ui.buildComponents(asts, root)
	if err != nil {
		return nil, err
	}

	return ui, nil
}

func (u *Ui) Update() error {
	if u.next != nil {
		u.nextTrans += SlideSpeed
		if u.nextTrans > 1 {
			u.current = u.next
			u.next = nil
			u.nextTrans = 0
		}
		return nil
	}

	fw, bw, s := u.controls.Update()

	if fw {
		tr := u.components[u.current].CurrentTarget()
		if len(tr.Children) > 0 {
			u.next = tr
			u.transDir = 1
		}
		return nil

	}
	if bw && u.current.Parent != nil {
		u.next = u.current.Parent
		u.transDir = -1
		return nil
	}

	u.components[u.current].Scroll(s)
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
}

func (u *Ui) Layout(_, _ int) (int, int) {
	return u.w, u.h
}

func (u *Ui) buildComponents(asts *assets.Loader, node *preset.Node) error {
	if len(node.Children) > 0 {
		list, err := NewList(asts, node)
		if err != nil {
			return err
		}
		u.components[node] = list
		for _, child := range node.Children {
			err := u.buildComponents(asts, child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
