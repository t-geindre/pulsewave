package ui

type Controls interface {
	Update() (horDelta, vertDelta int)
}

type MultiControls struct {
	controls []Controls
}

func NewMultiControls(ctrls ...Controls) *MultiControls {
	return &MultiControls{
		controls: ctrls,
	}
}

func (c *MultiControls) Update() (int, int) {
	horDelta, vertDelta := 0, 0
	for _, ctrl := range c.controls {
		hd, vd := ctrl.Update()
		horDelta += hd
		vertDelta += vd
	}
	return horDelta, vertDelta
}
