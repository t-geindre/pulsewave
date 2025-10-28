package dsp

type OscShape int

const (
	ShapeSine OscShape = iota
	ShapeSaw
	ShapeTriangle
	ShapeSquare
	ShapeNoise
)

type ShapeRegistry struct {
	shapes map[int]OscShape
}

func NewShapeRegistry() *ShapeRegistry {
	return &ShapeRegistry{}
}

func (s *ShapeRegistry) Set(id int, shape OscShape) {
	if s.shapes == nil {
		s.shapes = make(map[int]OscShape)
	}
	s.shapes[id] = shape
}

func (s *ShapeRegistry) Get(id int) OscShape {
	if shape, ok := s.shapes[id]; ok {
		return shape
	}

	return ShapeSine
}
