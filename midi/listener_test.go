package midi

import (
	"synth/msg"
	"testing"

	"github.com/rs/zerolog"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestListener_StopsGoroutines(t *testing.T) {
	defer goleak.VerifyNone(t)

	out := msg.NewQueue(1024)
	logger := zerolog.Nop()

	l := NewListener(logger, out)

	go l.ListenAll()
	l.Close()
}
