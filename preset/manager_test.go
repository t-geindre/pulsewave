package preset

import (
	"synth/dsp"
	"synth/msg"
	"testing"

	"github.com/rs/zerolog"
)

func TestManager_ProcessNoAlloc(t *testing.T) {
	// Send a save message
	inQueue := msg.NewQueue(1)
	inQueue.TryWrite(msg.Message{
		Kind: LoadSavePresetKind,
		Key:  0, // Default
		ValF: 1, // Save
	})

	messenger := msg.NewMessenger(inQueue, msg.NewQueue(1), 0)
	manager := NewManager(44100, zerolog.Nop(), messenger, "/dev/null")

	for i := 0; i < 16; i++ {
		manager.NoteOn(10+i, 1.0)
	}
	var block dsp.Block

	allocs := testing.AllocsPerRun(1000, func() {
		messenger.Process()
		manager.Process(&block)
	})
	if allocs != 0 {
		t.Errorf("expected 0 allocations, got %v", allocs)
	}
}
