package dsp

type OscShape int

const (
	ShapeSaw OscShape = iota
	ShapeTriangle
	ShapeSquare
	ShapeTableWave
)

type ShapeRegistry struct {
	shapes []OscShape
	tables []*Wavetable
}

func NewShapeRegistry() *ShapeRegistry {
	return &ShapeRegistry{
		shapes: make([]OscShape, 0),
		tables: make([]*Wavetable, 0),
	}
}

func (s *ShapeRegistry) Add(shape OscShape, table ...*Wavetable) float32 {
	var wt *Wavetable
	if len(table) > 0 {
		wt = table[0]
	} else {
		// avoid constant nil checking
		wt = NewZeroWavetable(1)
	}

	s.shapes = append(s.shapes, shape)
	s.tables = append(s.tables, wt)

	return float32(len(s.shapes) - 1)
}

func (s *ShapeRegistry) Set(idf float32, shape OscShape, table ...*Wavetable) {
	id := int(idf)
	var wt *Wavetable
	if len(table) > 0 {
		s.tables[id] = table[0]
	}
	s.shapes[id] = shape
	s.tables[id] = wt
}

func (s *ShapeRegistry) Get(idf float32) (OscShape, *Wavetable) {
	id := int(idf)
	return s.shapes[id], s.tables[id]
}
