package tree

import (
	"fmt"
	"synth/msg"
)

type SliderNode interface {
	ValueNode

	Step() float32
	Display() string
}

type sliderNode struct {
	min, max, step float32
	format         func(float32) string
	publish        func(uint8, float32)

	ValueNode
}

func NewSliderNode(label string, kind msg.Kind, key uint8, min, max, step float32, format func(float32) string) SliderNode {
	return &sliderNode{
		min:    min,
		max:    max,
		step:   step,
		format: format,

		ValueNode: NewValueNode(label, kind, key),
	}
}

func (n *sliderNode) Step() float32 {
	return n.step
}

func (n *sliderNode) SetVal(v float32) {
	if v < n.min {
		v = n.min
	}

	if v > n.max {
		v = n.max
	}

	n.ValueNode.SetVal(v)
}

func (n *sliderNode) Display() string {
	if n.format == nil {
		return fmt.Sprintf("%.2f", n.Val())
	}
	return n.format(n.Val())
}
