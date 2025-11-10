package preset

import (
	"synth/dsp"
	"synth/msg"
	"testing"
)

func TestPolysynth_ProcessNoAlloc(t *testing.T) {
	// Poll all parameters
	qIn, qOut := msg.NewQueue(64), msg.NewQueue(64)
	qIn.TryWrite(msg.Message{
		Kind: ParamPullAllKind,
	})

	synth := NewPolysynth(44100, qIn, qOut)
	for i := 0; i < 16; i++ {
		synth.voice.NoteOn(10+i, 1.0)
	}
	var block dsp.Block

	allocs := testing.AllocsPerRun(1000, func() {
		synth.Process(&block)
	})
	if allocs != 0 {
		t.Errorf("expected 0 allocations, got %v", allocs)
	}
}
