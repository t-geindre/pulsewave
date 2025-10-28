package audio

import (
	"encoding/binary"
	"math"
)

const (
	channels       = 2
	bytesPerSample = 4 // float32
	frameBytes     = channels * bytesPerSample
)

type Stream struct {
	source Source
	block  *Block
}

func NewStream(src Source) *Stream {
	return &Stream{
		source: src,
		block: &Block{
			L: [BlockSize]float32{},
			R: [BlockSize]float32{},
		},
	}
}

func (s *Stream) Read(p []byte) (int, error) {
	frames := len(p) / frameBytes
	if frames == 0 {
		return 0, nil
	}

	done := 0
	for done < frames {
		if s.block.left == 0 {
			s.pullBlock()
		}

		toCopy := s.block.left
		remain := frames - done
		if toCopy > remain {
			toCopy = remain
		}

		start := BlockSize - s.block.left

		for i := 0; i < toCopy; i++ {
			j := start + i
			off := (done + i) * frameBytes
			binary.LittleEndian.PutUint32(p[off+0:], math.Float32bits(s.block.L[j]))
			binary.LittleEndian.PutUint32(p[off+4:], math.Float32bits(s.block.R[j]))
		}

		done += toCopy
		s.block.left -= toCopy
	}

	return done * frameBytes, nil
}

func (s *Stream) pullBlock() {
	s.block.Cycle++
	s.block.left = BlockSize
	s.source.Process(s.block)
}
