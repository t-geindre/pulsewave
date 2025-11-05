package preset

type SelectorNode struct {
	options []*SelectorOption

	*ListNode
	*ParamNode
}

func NewSelectorNode(label string, key uint8, options ...*SelectorOption) *SelectorNode {
	return &SelectorNode{
		options:   options,
		ListNode:  NewListNode(label),
		ParamNode: NewParamNode(key),
	}
}

func (s *SelectorNode) Options() []*SelectorOption {
	return s.options
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
