package tree

type SelectorNode interface {
	Node
	ValueNode

	Options() []*SelectorOption
	RequiresValidation() bool
	Validate()
}

type selectorNode struct {
	options []*SelectorOption

	ValueNode
}

func NewSelectorNode(label string, key uint8, options ...*SelectorOption) SelectorNode {
	return &selectorNode{
		options:   options,
		ValueNode: NewValueNode(label, key),
	}
}

func (s *selectorNode) Options() []*SelectorOption {
	return s.options
}

func (s *selectorNode) RequiresValidation() bool {
	return false
}

func (s *selectorNode) Validate() {
	// no-op
}

type SelectorOption struct {
	icon  string
	label string
	value float32
}

func NewSelectorOption(label, icon string, value float32) *SelectorOption {
	return &SelectorOption{
		icon:  icon,
		label: label,
		value: value,
	}
}

func (o *SelectorOption) Icon() string {
	return o.icon
}

func (o *SelectorOption) Label() string {
	return o.label
}

func (o *SelectorOption) Value() float32 {
	return o.value
}
