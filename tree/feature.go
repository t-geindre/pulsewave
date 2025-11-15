package tree

const (
	FeatureOscilloscope = iota
)

type FeatureNode interface {
	Node
	Feature() int
}

type featureNode struct {
	Node
	feature int
}

func NewFeatureNode(label string, feature int) FeatureNode {
	return &featureNode{
		Node:    NewNode(label),
		feature: feature,
	}
}

func (f *featureNode) Feature() int {
	return f.feature
}
