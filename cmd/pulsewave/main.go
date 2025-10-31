package main

import (
	"os"
	"synth/assets"
	"synth/dsp"
	"synth/midi"
	"synth/msg"
	"synth/preset"
	"synth/ui"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/rs/zerolog"
)

func main() {
	const SampleRate = 44100

	// Debug mode
	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		debugMode = true
	}

	// Signal Chain
	synth := preset.NewPolysynth(SampleRate)
	headroom := dsp.NewVca(synth, dsp.NewParam(0.8))
	clean := dsp.NewLowPassSVF(SampleRate, headroom, dsp.NewParam(18000), dsp.NewParam(0.5))

	// Messaging
	router := msg.NewRouter(logger())
	midiInQ, midiInId := router.AddInput(1024)
	midiAudioOutQ, midiAudioOutId := router.AddOutput(1024)
	midiUiOutQ, miniUiOutId := router.AddOutput(1024)

	router.AddRoute(midiInId, midiAudioOutId, midi.MidiSource, midi.NoteOnKind)
	router.AddRoute(midiInId, midiAudioOutId, midi.MidiSource, midi.NoteOffKind)

	router.AddRoute(midiInId, miniUiOutId, midi.MidiSource, midi.ControlChangeKind)

	go router.Route()

	// Add midi player to synth
	midiPlayer := midi.NewPlayer(clean, synth, midiAudioOutQ)

	// Midi setup
	midi := midi.NewListener(logger())
	defer midi.Close()

	device, err := midi.FindDevice()
	onError(err, "failed to find device")

	err = midi.Listen(device, midiInQ)
	onError(err, "failed to listen to device")

	// Player
	ctx := audio.NewContext(SampleRate)
	player, err := ctx.NewPlayerF32(dsp.NewStream(midiPlayer))
	onError(err, "failed to create player")

	player.SetBufferSize(time.Millisecond * 20)
	player.Play()

	// Build menu
	menu := preset.NewMenu()

	// UI
	asts, err := assets.NewFromJson("assets/assets.json")
	onError(err, "failed to create assets loader")

	err = asts.Load()
	onError(err, "failed to load assets")

	ctrl := ui.NewControls(midiUiOutQ)

	gui, err := ui.NewUi(asts, ctrl, menu)
	onError(err, "failed to create gui")

	err = ebiten.RunGame(gui)
	onError(err, "failed to run gui")
}

func onError(err error, msg string) {
	l := logger()
	if err != nil {
		if debugMode {
			l.Panic().Err(err).Msg(msg)
		} else {
			l.Fatal().Err(err).Msg(msg)
		}
	}
}

var debugMode = false
var loggerInst *zerolog.Logger

func logger() zerolog.Logger {
	if loggerInst != nil {
		return *loggerInst
	}

	l := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.TimeOnly, NoColor: false},
	).With().
		Str("component", "main").
		Timestamp().
		Logger().
		Level(zerolog.InfoLevel)

	if debugMode {
		l = l.Level(zerolog.TraceLevel)
	}

	loggerInst = &l

	return l
}
