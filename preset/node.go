package preset

type Node interface {
	Label() string
	Children() []Node
	Parent() Node
	SetParent(Node)
	Append(Node)
	Prepend(Node)
	SetRoot(Node)
	Root() Node
}

type AttachFunc func(uint8, float32)

type ValueNode interface {
	Val() float32
	SetVal(float32)
	Key() uint8
	Attach(AttachFunc)
}

type OptionNode interface {
	Node
	ValueNode
	Options() []*SelectorOption
}

type OptionValidateNode interface {
	Validate()
}
