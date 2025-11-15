package dsp

type Callback struct {
	fn    func(block *Block)
	inner Node
}

func NewCallback(fn func(block *Block), inner Node) *Callback {
	return &Callback{
		fn:    fn,
		inner: inner,
	}
}

func (c *Callback) Process(block *Block) {
	c.fn(block)
	c.inner.Process(block)
}

func (c *Callback) Reset(soft bool) {
	c.inner.Reset(soft)
}
