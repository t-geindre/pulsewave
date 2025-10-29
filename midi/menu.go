package midi

type Node struct {
	Parent   *Node
	Label    string
	Bind     func(v int)
	Children []*Node
}

type Menu struct {
	wheel   *Wheel
	root    *Node
	current *Node
	cursor  int
}

func NewMenu() *Menu {
	m := &Menu{
		wheel: NewWheel(),
	}

	m.build()
	m.current = m.root

	return m
}

func (m *Menu) Wheel(v uint8) {

	if m.current.Bind == nil {
		m.cursor += m.wheel.Update(v) / 10
		for m.cursor < 0 {
			m.cursor += len(m.current.Children)
		}
		m.cursor = m.cursor % len(m.current.Children)
	}
}

func (m *Menu) Forward() {
}

func (m *Menu) Backward() {
}

func (m *Menu) build() {
	m.root = &Node{
		Label: "Menu",
		Children: []*Node{
			{
				Label: "Oscillator",
				Children: []*Node{
					{
						Label: "Waveform",
					},
					{
						Label: "Frequency",
					},
				},
			},
			{
				Label: "Filter",
				Children: []*Node{
					{
						Label: "Cutoff",
					},
					{
						Label: "Resonance",
					},
				},
			},
		},
	}
}
