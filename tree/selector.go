package tree

import "synth/msg"

type SelectorNode interface {
	Node
	ValueNode

	Options() []*SelectorOption
	RequiresValidation() bool
	Validate()
}

type selectorNode struct {
	options            []*SelectorOption
	requiresValidation bool
	val                float32

	ValueNode
}

func NewSelectorNode(label string, kind msg.Kind, key uint8, options ...*SelectorOption) SelectorNode {
	return &selectorNode{
		options:   options,
		ValueNode: NewValueNode(label, kind, key),
	}
}

func NewValidatingSelectorNode(label string, kind msg.Kind, key uint8, options ...*SelectorOption) SelectorNode {
	return &selectorNode{
		options:            options,
		requiresValidation: true,
		ValueNode:          NewValueNode(label, kind, key),
	}
}

func (s *selectorNode) Options() []*SelectorOption {
	return s.options
}

func (s *selectorNode) RequiresValidation() bool {
	return s.requiresValidation
}

func (s *selectorNode) Validate() {
	if s.requiresValidation {
		s.ValueNode.SetValAndPublish(s.val)
	}
}

func (s *selectorNode) SetVal(f float32) {
	if s.requiresValidation {
		s.val = f
	} else {
		s.ValueNode.SetVal(f)
	}
}

func (s *selectorNode) Val() float32 {
	if s.requiresValidation {
		return s.val
	}
	return s.ValueNode.Val()
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
