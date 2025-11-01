package dsp

type OscShape int

const (
	ShapeSine OscShape = iota
	ShapeSaw
	ShapeTriangle
	ShapeSquare
	ShapeNoise
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

func (s *ShapeRegistry) Add(shape OscShape, table ...*Wavetable) int {
	var wt *Wavetable
	if len(table) > 0 {
		wt = table[0]
	} else {
		// avoid constant nil checking
		wt = NewZeroWavetable(1)
	}

	s.shapes = append(s.shapes, shape)
	s.tables = append(s.tables, wt)

	return len(s.shapes) - 1
}

func (s *ShapeRegistry) Set(id int, shape OscShape, table ...*Wavetable) {
	var wt *Wavetable
	if len(table) > 0 {
		s.tables[id] = table[0]
	}
	s.shapes[id] = shape
	s.tables[id] = wt
}

func (s *ShapeRegistry) Get(id int) (OscShape, *Wavetable) {
	return s.shapes[id], s.tables[id]
}
