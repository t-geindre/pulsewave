package audio

const BlockSize = 512

type Block struct {
	// Cycle Unique identifier for the current block
	Cycle uint64

	// Channels L/R
	L, R [BlockSize]float32

	// Remaining samples to process
	left int
}
