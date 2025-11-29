package tree

type RedirectionNode interface {
	Node
	GetRedirection() Node
}

type targetNode struct {
	Node
}

func NewRedirectionNode(label string) RedirectionNode {
	return &targetNode{
		Node: NewNode(label),
	}
}

func (t *targetNode) GetRedirection() Node {
	return t.Context().EnterSubTree(
		t.Parent(),
		t.Root(),
	)
}
