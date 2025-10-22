package audio

import "math"

type Stream struct {
	source Source
}

func NewStream(source Source) *Stream {
	return &Stream{source: source}
}

func (s *Stream) Read(p []byte) (int, error) {
	frames := len(p) / 8 // 2 chan * 4 octets/sample
	for i := 0; i < frames; i++ {
		vl, vr := s.source.NextValue()
		off := i * 8

		// Left
		bitsL := math.Float32bits(float32(vl))
		p[off+0] = byte(bitsL)
		p[off+1] = byte(bitsL >> 8)
		p[off+2] = byte(bitsL >> 16)
		p[off+3] = byte(bitsL >> 24)

		// Right
		bitsR := math.Float32bits(float32(vr))
		p[off+4] = byte(bitsR)
		p[off+5] = byte(bitsR >> 8)
		p[off+6] = byte(bitsR >> 16)
		p[off+7] = byte(bitsR >> 24)
	}

	return frames * 8, nil
}
