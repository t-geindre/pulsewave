package midi

import (
	"sync"
	"synth/msg"
	"time"

	"github.com/rs/zerolog"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type device struct {
	name  string
	stop  func()
	found bool
}

type Listener struct {
	logger  zerolog.Logger
	out     *msg.Queue
	devices []*device
	msgs    chan msg.Message
	done    chan struct{}
	once    sync.Once
}

func NewListener(log zerolog.Logger, out *msg.Queue) *Listener {
	return &Listener{
		logger: log.With().Str("component", "MIDI listener").Logger(),
		out:    out,
		msgs:   make(chan msg.Message, 1024),
		done:   make(chan struct{}),
	}
}

func (l *Listener) ListenAll() {
	// Message writer
	go func() {
		for {
			select {
			case <-l.done:
				return
			case m := <-l.msgs:
				l.out.TryWrite(m)
			}
		}
	}()

	// Device polling
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-l.done:
			// Done, stop all devices
			for _, dev := range l.devices {
				dev.stop()
			}
			return
		case <-ticker.C:
			// Mark all devices as not found
			for _, dev := range l.devices {
				dev.found = false
			}

			// Check for new devices
			for _, in := range midi.GetInPorts() {
				dev := l.findDevice(in)
				if dev == nil {
					l.listenDevice(in)
					continue
				}
				dev.found = true
			}

			// Remove devices that are no longer connected
			var activeDevices []*device
			for _, dev := range l.devices {
				if !dev.found {
					dev.stop()
				} else {
					activeDevices = append(activeDevices, dev)
				}
			}
			l.devices = activeDevices
		}
	}
}

// Close stops the listener and all associated goroutines.
func (l *Listener) Close() {
	l.once.Do(func() {
		close(l.done)
	})
}

func (l *Listener) findDevice(in drivers.In) *device {
	for _, dev := range l.devices {
		if dev.name == in.String() {
			return dev
		}
	}
	return nil
}

func (l *Listener) listenDevice(in drivers.In) {
	stop, err := midi.ListenTo(
		in,
		l.handleMessage,
		midi.HandleError(func(err error) {
			l.logger.Error().Err(err).Msg("listen error")
		}),
		midi.UseSysEx(),
	)

	if err != nil {
		l.logger.Error().Err(err).Str("device", in.String()).Msg("failed to listen")
		return
	}

	l.logger.Info().Str("device", in.String()).Msg("listening")

	l.devices = append(l.devices, &device{
		name: in.String(),
		stop: func() {
			stop()
			l.logger.Info().Str("device", in.String()).Msg("stopped listening")
		},
		found: true,
	})
}

func (l *Listener) handleMessage(message midi.Message, _ int32) {
	var ch, key, val8 uint8
	var val16 int16
	switch {
	case message.GetNoteStart(&ch, &key, &val8):
		l.send(msg.Message{
			Kind: NoteOnKind,
			Key:  key,
			Val8: val8,
			Chan: ch,
		})
		l.logger.Debug().Uint8("channel", ch).Uint8("key", key).Uint8("val8", val8).Msg("Note ON")

	case message.GetNoteEnd(&ch, &key):
		l.send(msg.Message{
			Kind: NoteOffKind,
			Key:  key,
			Val8: val8,
			Chan: ch,
		})
		l.logger.Debug().Uint8("channel", ch).Uint8("key", key).Uint8("val8", val8).Msg("Note OFF")
	case message.GetControlChange(&ch, &key, &val8):
		l.send(msg.Message{
			Kind: ControlChangeKind,
			Key:  key,
			Val8: val8,
			Chan: ch,
		})
		l.logger.Debug().Uint8("channel", ch).Uint8("controller", key).Uint8("value", val8).Msg("Control Change")
	case message.GetPitchBend(&ch, &val16, nil):
		l.send(msg.Message{
			Kind:  PitchBendKind,
			Val16: val16,
			Chan:  ch,
		})
		l.logger.Debug().Uint8("channel", ch).Int16("value", val16).Msg("Pitch Bend")
	default:
		l.logger.Debug().Str("msg", message.String()).Msg("unknown message")
	}
}

func (l *Listener) send(m msg.Message) {
	select {
	case <-l.done:
		// closing, drop message
		return
	case l.msgs <- m:
		return
	}
}
