package dsp

type Shape float32

const (
	Sine Shape = iota
	Saw
	Triangle
	Square
	Noise
)

type ShapeRegistry struct {
	shapes map[int]Shape
}

func NewShapeRegistry() *ShapeRegistry {
	return &ShapeRegistry{}
}

func (s *ShapeRegistry) Set(id int, shape Shape) {
	if s.shapes == nil {
		s.shapes = make(map[int]Shape)
	}
	s.shapes[id] = shape
}

func (s *ShapeRegistry) Get(id int) Shape {
	if shape, ok := s.shapes[id]; ok {
		return shape
	}

	return Sine
}
