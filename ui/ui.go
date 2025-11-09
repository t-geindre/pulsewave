package ui

import (
	"fmt"
	"os"
	"synth/assets"
	"synth/msg"
	"synth/preset"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"google.golang.org/protobuf/proto"
)

// Todo get it fron config
const (
	BodyWidth  = 376
	BodyHeight = 211
	BodyStartX = 51
	BodyStartY = 70
	SlideSpeed = .08
)

type Ui struct {
	background *ebiten.Image
	w, h       int
	controls   Controls

	messenger *Messenger

	components map[preset.Node]Component
	current    preset.Node
	next       preset.Node
	nextTrans  float64
	transDir   float64
	tree       *preset.Tree

	transLeft    *ebiten.Image
	transRight   *ebiten.Image
	bodyClipMask *ebiten.Image
}

func NewUi(asts *assets.Loader, inQ, outQ *msg.Queue) (*Ui, error) {
	// BG + window size accordingly
	bg, err := asts.GetImage("ui/background")
	if err != nil {
		return nil, err
	}
	bds := bg.Bounds()

	ebiten.SetWindowSize(bds.Dx(), bds.Dy())
	ebiten.SetWindowTitle("Pulsewave")

	// Menu tree + controls
	midiCtrls := NewMidiControls()
	tree := preset.NewTree()

	messenger := NewMessenger(tree, midiCtrls, inQ, outQ)

	ctrls := NewMultiControls(
		NewPlayControls(messenger),
		NewKeyboardControls(),
		midiCtrls,
	)

	ui := &Ui{
		background:   bg,
		w:            bds.Dx(),
		h:            bds.Dy(),
		components:   make(map[preset.Node]Component),
		messenger:    messenger,
		current:      tree.Node,
		tree:         tree,
		controls:     ctrls,
		bodyClipMask: ebiten.NewImage(BodyWidth, BodyHeight),
		transLeft:    ebiten.NewImage(BodyWidth, BodyHeight),
		transRight:   ebiten.NewImage(BodyWidth, BodyHeight),
	}

	err = ui.buildComponents(asts, ui.current)
	if err != nil {
		return nil, err
	}

	// todo Move elsewhere
	presetData, err := os.ReadFile("assets/presets/saved.preset")
	if err != nil {
		return nil, err
	}
	var pp preset.ProtoPreset
	err = proto.Unmarshal(presetData, &pp)
	if err != nil {
		return nil, err
	}
	p := preset.NewFromProto(&pp)
	ui.tree.LoadPreset(p)
	fmt.Printf("loaded %f\n", p[preset.Osc0Shape].GetBase())

	return ui, nil
}

func (u *Ui) Update() error {
	// Todo implement a proper save dialog
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		p := u.tree.GetPreset()
		pp := p.ToProto()
		raw, err := proto.Marshal(pp)
		if err != nil {
			return err
		}
		err = os.WriteFile("assets/presets/saved.preset", raw, 0644)
		if err != nil {
			return err
		}
		fmt.Printf("saved %f\n", p[preset.Osc0Shape].GetBase())
	}
	u.messenger.Update()

	if u.next != nil {
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
}

func (u *Ui) Layout(_, _ int) (int, int) {
	return u.w, u.h
}

func (u *Ui) buildComponents(asts *assets.Loader, n preset.Node) error {
	switch node := n.(type) {
	case *preset.SliderNode:
		comp, err := NewSlider(asts, node)
		if err != nil {
			return err
		}
		u.components[node] = comp
	case *preset.ListNode:
		if len(node.Children()) > 0 {
			comp, err := NewList(asts, node)
			if err != nil {
				return err
			}
			u.components[node] = comp
		}
	case *preset.SelectorNode:
		comp, err := NewSelector(asts, node)
		if err != nil {
			return err
		}
		u.components[node] = comp
	}

	for _, child := range n.Children() {
		err := u.buildComponents(asts, child)
		if err != nil {
			return err
		}
	}

	return nil
}
