package audio

type TrackSet struct {
	sources []Source
	gains   []float64
	loop    bool
}

func NewTrackSet(sources ...Source) *TrackSet {
	t := &TrackSet{}
	for _, src := range sources {
		t.Append(src, 1)
	}
	return t
}

func (m *TrackSet) NextValue() (float64, float64) {
	vl, vr := .0, .0
	active := false

	for i, src := range m.sources {
		if src.IsActive() {
			active = true
		}
		l, r := src.NextValue()
		vl += l * m.gains[i]
		vr += r * m.gains[i]
	}

	if !active && m.loop {
		m.Reset()
	}

	return vl, vr
}

func (m *TrackSet) Append(source Source, gain float64) {
	m.sources = append(m.sources, source)
	m.gains = append(m.gains, gain)
}

func (m *TrackSet) SetLoop(loop bool) {
	m.loop = loop
}

func (m *TrackSet) IsActive() bool {
	for _, src := range m.sources {
		if src.IsActive() {
			return true
		}
	}
	return false
}

func (m *TrackSet) Reset() {
	for _, src := range m.sources {
		src.Reset()
	}
}

func (m *TrackSet) NoteOn(_, _ float64) {
}

func (m *TrackSet) NoteOff(float64) {
}
