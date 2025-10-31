package midi

import (
	"errors"
	"strings"
	"synth/msg"

	"github.com/rs/zerolog"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

var ErrNoMidiDevice = errors.New("no midi device found")

type Listener struct {
	logger zerolog.Logger
	stop   func()
}

func NewListener(log zerolog.Logger) *Listener {
	return &Listener{
		logger: log.With().Str("component", "MIDI listener").Logger(),
	}
}

func (l *Listener) FindDevice() (drivers.In, error) {
	ips := midi.GetInPorts()
	devices := make([]drivers.In, 0)

	for _, port := range ips {
		if strings.Contains(port.String(), "Midi Through") { // Linux virtual port, skip
			continue
		}
		l.logger.Info().Str("device", port.String()).Msg("device found")
		devices = append(devices, port)
	}

	if len(devices) == 0 {
		return nil, ErrNoMidiDevice
	}

	if len(devices) > 1 {
		l.logger.Warn().Str("device", devices[0].String()).Msg("multiple devices, using first")
	}

	return devices[0], nil
}

func (l *Listener) Listen(device drivers.In, queue *msg.Queue) error {
	l.Stop()

	var err error
	l.stop, err = midi.ListenTo(device, func(message midi.Message, _ int32) {
		var ch, key, val8 uint8
		var val16 int16
		switch {
		case message.GetNoteStart(&ch, &key, &val8):
			queue.TryWrite(msg.Message{
				Source: MidiSource,
				Kind:   NoteOnKind,
				Key:    key,
				Val8:   val8,
				Chan:   ch,
			})
			l.logger.Debug().Uint8("channel", ch).Uint8("key", key).Uint8("val8", val8).Msg("Note ON")

		case message.GetNoteEnd(&ch, &key):
			queue.TryWrite(msg.Message{
				Source: MidiSource,
				Kind:   NoteOffKind,
				Key:    key,
				Val8:   val8,
			})
			l.logger.Debug().Uint8("channel", ch).Uint8("key", key).Uint8("val8", val8).Msg("Note OFF")
		case message.GetControlChange(&ch, &key, &val8):
			queue.TryWrite(msg.Message{
				Source: MidiSource,
				Kind:   ControlChangeKind,
				Key:    key,
				Val8:   val8,
				Chan:   ch,
			})
			l.logger.Debug().Uint8("channel", ch).Uint8("controller", key).Uint8("value", val8).Msg("Control Change")
		case message.GetPitchBend(&ch, &val16, nil):
			queue.TryWrite(msg.Message{
				Source: MidiSource,
				Kind:   PitchBendKind,
				Val16:  val16,
				Chan:   ch,
			})
			l.logger.Debug().Uint8("channel", ch).Int16("value", val16).Msg("Pitch Bend")
		default:
			l.logger.Debug().Str("msg", message.String()).Msg("message ignored")
		}
	})

	return err
}

func (l *Listener) Stop() {
	if l.stop != nil {
		l.stop()
	}
}

func (l *Listener) Close() {
	l.Stop()
	midi.CloseDriver()
}
