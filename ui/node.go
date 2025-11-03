package ui

type Node interface {
	Label() string
	Children() []Node
	Parent() Node
	SetParent(Node)
	Append(Node)
}

type ListNode struct {
	label    string
	children []Node
	parent   Node
}

func NewListNode(label string, children ...Node) *ListNode {
	n := &ListNode{
		label:    label,
		children: []Node{},
	}
	for _, c := range children {
		n.Append(c)
	}
	return n
}

func (n *ListNode) Label() string {
	return n.label
}

func (n *ListNode) Children() []Node {
	return n.children
}

func (n *ListNode) Parent() Node {
	return n.parent
}

func (n *ListNode) SetParent(p Node) {
	n.parent = p
}

func (n *ListNode) Append(child Node) {
	child.SetParent(n)
	n.children = append(n.children, child)
}

type SliderNode struct {
	Unit           string
	Min, Max, Step uint16
	Key            uint8

	*ListNode
}

func NewSliderNode(label, unit string, key uint8, min, max, step uint16) *SliderNode {
	return &SliderNode{
		Unit: unit,
		Key:  key,
		Min:  min,
		Max:  max,
		Step: step,

		ListNode: NewListNode(label),
	}
}
