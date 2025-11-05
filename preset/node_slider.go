package preset

import "fmt"

type SliderNode struct {
	min, max, step float32
	format         func(float32) string
	publish        func(uint8, float32)

	*ListNode
	*ParamNode
}

func NewSliderNode(label string, key uint8, min, max, step float32, format func(float32) string) *SliderNode {
	return &SliderNode{
		min:    min,
		max:    max,
		step:   step,
		format: format,

		ParamNode: NewParamNode(key),
		ListNode:  NewListNode(label),
	}
}

func (n *SliderNode) Step() float32 {
	return n.step
}

func (n *SliderNode) SetVal(v float32) {
	if v < n.min {
		v = n.min
	}

	if v > n.max {
		v = n.max
	}

	n.ParamNode.SetVal(v)
}

func (n *SliderNode) Display() string {
	if n.format == nil {
		return fmt.Sprintf("%.2f", n.val)
	}
	return n.format(n.val)
}
