package audio

type TrackSet struct {
	sources []Source
}

func NewTrackSet(sources ...Source) *TrackSet {
	return &TrackSet{sources: sources}
}

func (m *TrackSet) NextSample() float64 {
	v := 0.0

	for _, src := range m.sources {
		v += src.NextSample()
	}

	return v
}

func (m *TrackSet) Append(sources ...Source) {
	m.sources = append(m.sources, sources...)
}
