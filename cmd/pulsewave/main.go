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
	if len(os.Args) > 1 && os.Args[1] == "--debugMode" {
		debugMode = true
	}

	// Signal Chain
	synth := preset.NewPolysynth(SampleRate)
	headroom := dsp.NewVca(synth, dsp.NewParam(0.8))
	clean := dsp.NewLowPassSVF(SampleRate, headroom, dsp.NewParam(18000), dsp.NewParam(0.5))

	// Messaging
	router := msg.NewRouter(logger())
	midiInQ, midiInId := router.AddInput(1024)
	midiOutQ, midiOutId := router.AddOutput(1024)

	router.AddRoute(midiInId, midiOutId, midi.MidiSource, midi.NoteOnKind)
	router.AddRoute(midiInId, midiOutId, midi.MidiSource, midi.NoteOffKind)

	go router.Route()

	// Add midi player to synth
	midiPlayer := midi.NewPlayer(clean, synth, midiOutQ)

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

	// UI
	for {
		time.Sleep(time.Millisecond * 100)
	}
	asts := assets.NewLoader()
	asts.AddImage("ui/background", "assets/imgs/background.png")
	asts.AddImage("ui/arrow", "assets/imgs/arrow.png")
	asts.AddImage("ui/selected", "assets/imgs/selected.png")
	asts.AddFont("ui/font", "assets/fonts/Roboto-Medium.ttf") // Roboto medium, letter spacing 3, size 20, color white
	asts.AddFace("ui/face", "ui/font", 21)
	asts.MustLoad()

	err = ebiten.RunGame(ui.NewUi(asts))
	onError(err, "failed to run ui")
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
