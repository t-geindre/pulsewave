package tree

type Context struct {
	subTreeFrom Node
	subTreeTo   Node
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) EnterSubTree(from, to Node) Node {
	c.subTreeFrom = from
	c.subTreeTo = to
	return to
}

func (c *Context) IsLeavingSubTree(current Node) Node {
	if c.subTreeFrom != nil && current == c.subTreeTo {
		target := c.subTreeFrom
		c.subTreeFrom = nil
		c.subTreeTo = nil
		return target
	}
	return nil
}
