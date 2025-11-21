package ui

import (
	"fmt"
	"synth/assets"
	"synth/tree"
)

var ErrorUnknownNodeType = fmt.Errorf("unknown node type or empty node")
var ErrorUnknownFeatureType = fmt.Errorf("unknown feature type")

type Components map[tree.Node]Component

func NewComponents(asts *assets.Loader, node tree.Node, audioQ *AudioQueue) (Components, error) {
	c := make(Components)

	err := c.build(asts, node, audioQ)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c Components) build(asts *assets.Loader, node tree.Node, aq *AudioQueue) error {
	comp, err := c.nodeComponent(asts, node, aq)
	if err != nil {
		return err
	}

	c[node] = comp

	for _, child := range node.Children() {
		err := c.build(asts, child, aq)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Components) nodeComponent(asts *assets.Loader, node tree.Node, aq *AudioQueue) (Component, error) {
	switch node := node.(type) {
	case tree.SliderNode:
		return NewSlider(asts, node)
	case tree.SelectorNode:
		return NewSelector(asts, node)
	case tree.FeatureNode:
		switch node.Feature() {
		case tree.FeatureOscilloscope:
			return NewOscilloscope(aq, 16384)
		default:

			return nil, ErrorUnknownFeatureType
		}
	default:
		if len(node.Children()) > 0 {
			return NewList(asts, node)
		}
		return nil, nil
	}
}
