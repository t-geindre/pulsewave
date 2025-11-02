package preset

type Node struct {
	Label    string
	Children []*Node
	Parent   *Node
}

func NewMenu() *Node {
	root := &Node{
		Children: []*Node{
			{
				Label: "Oscillators",
				Children: []*Node{
					{
						Label: "Oscillator 1",
						Children: []*Node{
							{Label: "Waveform"},
							{Label: "Octave"},
							{Label: "Semitone"},
							{Label: "Detune"},
							{Label: "Gain"},
						},
					},
					{
						Label: "Oscillator 2",
						Children: []*Node{
							{Label: "Waveform"},
							{Label: "Octave"},
							{Label: "Semitone"},
							{Label: "Detune"},
							{Label: "Gain"},
						},
					},
					{
						Label: "Oscillator 3",
						Children: []*Node{
							{Label: "Waveform"},
							{Label: "Octave"},
							{Label: "Semitone"},
							{Label: "Detune"},
							{Label: "Gain"},
						},
					},
				},
			},
			{
				Label: "Modulation",
				Children: []*Node{
					{Label: "Amplitude"},
					{Label: "Cutoff"},
					{Label: "Resonance"},
					{Label: "Pitch"},
				},
			},
			{
				Label: "Effects",
				Children: []*Node{
					{Label: "FB Delay"},
				},
			},
			{
				Label: "Filters",
				Children: []*Node{
					{Label: "Cutoff"},
				},
			},
			{
				Label: "Visualizer",
				Children: []*Node{
					{Label: "Spectrum"},
					{Label: "Oscilloscope"},
				},
			},
			{
				Label: "Presets",
				Children: []*Node{
					{Label: "Load preset"},
					{Label: "Save preset"},
					// Todo ON/OFF, we should have a dedicated entry for toggles with smth like a checkbox
					{Label: "Auto save"},
				},
			},
			{
				Label: "Settings",
				Children: []*Node{
					{Label: "General"},
					{Label: "Master gain"},
					{Label: "Controller"},
					{Label: "Pitch bend"},
					{Label: "About"},
				},
			},
		},
	}

	assignParents(root)

	return root
}

func assignParents(node *Node) {
	for _, child := range node.Children {
		child.Parent = node
		assignParents(child)
	}
}
