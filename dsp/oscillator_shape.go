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
	shapes map[int]OscShape
	tables map[int]*Wavetable
}

func NewShapeRegistry() *ShapeRegistry {
	return &ShapeRegistry{}
}

func (s *ShapeRegistry) Set(id int, shape OscShape, table ...*Wavetable) {
	if s.shapes == nil {
		s.shapes = make(map[int]OscShape)
	}

	if len(table) > 0 {
		if s.tables == nil {
			s.tables = make(map[int]*Wavetable)
		}
		s.tables[id] = table[0]
	}

	s.shapes[id] = shape
}

func (s *ShapeRegistry) Get(id int) (OscShape, *Wavetable) {
	if shape, ok := s.shapes[id]; ok {
		var table *Wavetable
		if s.tables != nil {
			table = s.tables[id]
		}
		return shape, table
	}

	return ShapeSine, nil
}
