package audio

type TrackSet struct {
	sources []Source
	gains   []float64
	loop    bool
}

func NewTrackSet(sources ...Source) *TrackSet {
	return &TrackSet{sources: sources}
}

func (m *TrackSet) NextSample() float64 {
	v := 0.0
	active := false

	for i, src := range m.sources {
		if acc, ok := src.(Resettable); ok {
			if acc.IsActive() {
				active = true
			}
		}
		v += src.NextSample() * m.gains[i]
	}

	if !active && m.loop {
		m.Reset()
	}

	return v
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
		if acc, ok := src.(Resettable); ok {
			if acc.IsActive() {
				return true
			}
		}
	}
	return false
}

func (m *TrackSet) Reset() {
	for _, src := range m.sources {
		if acc, ok := src.(Resettable); ok {
			acc.Reset()
		}
	}
}
