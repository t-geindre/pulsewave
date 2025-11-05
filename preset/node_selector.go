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
