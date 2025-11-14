package main

import (
	"flag"
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

	// Flags
	helpF := flag.Bool("help", false, "show help")
	debugF := flag.Bool("debug", false, "enable debug mode")
	fullsF := flag.Bool("full-screen", false, "enable full screen mode")
	buffF := flag.Int("buffer", 25, "disable gui")
	flag.Parse()

	// Help
	if *helpF {
		flag.Usage()
		os.Exit(0)
	}

	// Debug mode
	if *debugF {
		debugMode = true
	}

	// Full screen
	if *fullsF {
		ebiten.SetFullscreen(true)
	}

	// Messaging : Midi
	router := msg.NewRouter(logger())
	midiInQ := router.AddInput(1024)
	audioOutQ := router.AddOutput(1024)
	audioInQ := router.AddInput(1024)

	uiOutQ := router.AddOutput(1024)
	uiInQ := router.AddInput(1024)

	router.AddRoute(midiInQ, midi.NoteOnKind, audioOutQ)
	router.AddRoute(midiInQ, midi.NoteOffKind, audioOutQ)
	router.AddRoute(midiInQ, midi.PitchBendKind, audioOutQ)

	router.AddRoute(midiInQ, midi.ControlChangeKind, uiOutQ)

	router.AddRoute(uiInQ, preset.ParamUpdateKind, audioOutQ)
	router.AddRoute(uiInQ, preset.ParamPullAllKind, audioOutQ)
	router.AddRoute(uiInQ, midi.NoteOnKind, audioOutQ)
	router.AddRoute(uiInQ, midi.NoteOffKind, audioOutQ)

	router.AddRoute(audioInQ, preset.ParamUpdateKind, uiOutQ)

	go router.Route()

	// Main signal + audio tap	for UI
	synth := preset.NewPolysynth(SampleRate, audioOutQ, audioInQ)

	uiAudioQueue := ui.NewAudioQueue(32) // 32 blocks x 256 samples
	synthTap := ui.NewAudioPuller(synth, uiAudioQueue)

	// Clean signal
	headroom := dsp.NewVca(synthTap, dsp.NewParam(0.9))
	clean := dsp.NewLowPassSVF(SampleRate, headroom, dsp.NewParam(18000), dsp.NewParam(0.5))

	// Midi setup
	mdi := midi.NewListener(logger(), midiInQ)
	defer mdi.Close()
	go mdi.ListenAll()

	// Player
	ctx := audio.NewContext(SampleRate)
	player, err := ctx.NewPlayerF32(dsp.NewStream(clean))
	onError(err, "failed to create player")
	defer player.Close()

	player.SetBufferSize(time.Millisecond * time.Duration(*buffF))
	player.Play()

	// UI
	asts, err := assets.NewFromJson("assets/assets.json")
	onError(err, "failed to create assets loader")

	err = asts.Load()
	onError(err, "failed to load assets")

	gui, err := ui.NewUi(asts, uiOutQ, uiInQ, uiAudioQueue, logger(), debugMode)
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
