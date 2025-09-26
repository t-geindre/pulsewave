package audio

import "math"

type Stream struct {
	source Source
}

func NewStream(source Source) *Stream {
	return &Stream{source: source}
}

func (s *Stream) Read(p []byte) (int, error) {
	frames := len(p) / 8 // 2 canaux * 4 octets/échantillon
	for i := 0; i < frames; i++ {
		v := float32(s.source.NextSample())

		bits := math.Float32bits(v)
		off := i * 8

		// Left
		p[off+0] = byte(bits)
		p[off+1] = byte(bits >> 8)
		p[off+2] = byte(bits >> 16)
		p[off+3] = byte(bits >> 24)

		// Right (same, mono)
		p[off+4] = byte(bits)
		p[off+5] = byte(bits >> 8)
		p[off+6] = byte(bits >> 16)
		p[off+7] = byte(bits >> 24)
	}

	return frames * 8, nil
}
