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

type ParameterNode struct {
	min, max, step, val float32
	key                 uint8
	format              func(float32) string
	publish             func(uint8, float32)

	*ListNode
}

func NewParameterNode(label string, key uint8, min, max, step float32, format func(float32) string) *ParameterNode {
	return &ParameterNode{
		key:    key,
		min:    min,
		max:    max,
		step:   step,
		format: format,

		ListNode: NewListNode(label),
	}
}

func (n *ParameterNode) Key() uint8 {
	return n.key
}

func (n *ParameterNode) Min() float32 {
	return n.min
}

func (n *ParameterNode) Max() float32 {
	return n.max
}

func (n *ParameterNode) Val() float32 {
	return n.val
}

func (n *ParameterNode) Step() float32 {
	return n.step
}

func (n *ParameterNode) SetVal(v float32) {
	if n.val == v {
		return
	}

	if v < n.min {
		v = n.min
	}

	if v > n.max {
		v = n.max
	}

	if n.publish != nil {
		n.publish(n.key, v)
	}

	n.val = v
}

func (n *ParameterNode) Display() string {
	return n.format(n.val)
}

func (n *ParameterNode) Attach(publish func(uint8, float32)) {
	n.publish = publish
}
