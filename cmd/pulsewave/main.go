package main

import (
	"flag"
	"os"
	"synth/assets"
	"synth/dsp"
	"synth/midi"
	"synth/msg"
	"synth/preset"
	"synth/settings"
	"synth/tree"
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

	// Full screen + window tilte
	if *fullsF {
		ebiten.SetFullscreen(true)
	}
	ebiten.SetWindowTitle("Pulsewave")

	// Messaging
	router := msg.NewRouter(logger().With().Str("component", "router").Logger())
	midiInQ := router.AddInput(1024)

	audioOutQ := router.AddOutput(1024)
	audioInQ := router.AddInput(1024)

	uiOutQ := router.AddOutput(1024)
	uiInQ := router.AddInput(1024)

	setsInQ := router.AddInput(1024)
	setsOutQ := router.AddOutput(1024)

	// Routing: MIDI to audio
	router.AddRoute(midiInQ, midi.NoteOnKind, audioOutQ)
	router.AddRoute(midiInQ, midi.NoteOffKind, audioOutQ)
	router.AddRoute(midiInQ, midi.PitchBendKind, audioOutQ)

	// Routing: UI to audio
	router.AddRoute(uiInQ, preset.PresetLoadSaveKind, audioOutQ)
	router.AddRoute(uiInQ, preset.PresetUpdateKind, audioOutQ)
	router.AddRoute(uiInQ, midi.NoteOnKind, audioOutQ)
	router.AddRoute(uiInQ, midi.NoteOffKind, audioOutQ)

	// Routing: audio to UI
	router.AddRoute(audioInQ, preset.PresetUpdateKind, uiOutQ)

	// Routing: settings to audio/UI + ui to settings
	router.AddRoute(uiInQ, settings.SettingUpdateKind, setsOutQ)
	router.AddRoute(setsInQ, settings.SettingUpdateKind, audioOutQ)
	router.AddRoute(setsInQ, settings.SettingUpdateKind, uiOutQ)

	go router.Route()

	// Main signal
	audioMessenger := msg.NewMessenger(audioOutQ, audioInQ, 10)
	presetManager := preset.NewManager(
		SampleRate,
		logger().With().Str("component", "presets").Logger(),
		audioMessenger,
		"assets/presets",
	)

	// Audio messenger injection
	withMessenger := dsp.NewCallback(func(block *dsp.Block) {
		audioMessenger.Process()
	}, presetManager)

	audioMessenger.RegisterHandler(midi.NewPlayer(presetManager))
	audioMessenger.RegisterHandler(presetManager)

	// Audio tap
	uiAudioQueue := ui.NewAudioQueue(32) // 32 blocks x 256 samples
	synthTap := ui.NewAudioPuller(withMessenger, uiAudioQueue)

	// Clean presetManager
	clean := dsp.NewLowPassSVF(SampleRate, synthTap, dsp.NewParam(18000), dsp.NewParam(0.5))

	// Midi setup
	mdi := midi.NewListener(
		logger().With().Str("component", "midi").Logger(),
		midiInQ,
	)
	defer mdi.Close()
	go mdi.ListenAll()

	// Player
	ctx := audio.NewContext(SampleRate)
	player, err := ctx.NewPlayerF32(dsp.NewStream(clean))
	onError(err, "failed to create player")
	defer player.Close()

	player.SetBufferSize(time.Millisecond * time.Duration(*buffF))
	player.Play()

	// Settings
	sets := settings.NewSettings(
		"assets/settings.cfg",
		msg.NewMessenger(setsOutQ, setsInQ, 0),
		logger().With().Str("component", "settings").Logger(),
	)

	defer sets.Close()
	defer sets.Persist()

	// Assets
	asts, err := assets.NewFromJson("assets/assets.json")
	onError(err, "failed to create assets loader")

	err = asts.Load()
	onError(err, "failed to load assets")

	// UI Messenger
	uiMessenger := msg.NewMessenger(uiOutQ, uiInQ, 0)

	// Menu tree
	menuTree := tree.NewTree(presetManager.GetPresets())
	menuTree.Attach(uiMessenger)

	// UI Components
	components, err := ui.NewComponents(asts, menuTree, uiAudioQueue)
	onError(err, "failed to create ui components")

	// Controls
	controls := ui.NewMultiControls(
		ui.NewPlayControls(uiMessenger),
		ui.NewKeyboardControls(),
		ui.NewTouchControls(),
	)

	// GUI
	gui, err := ui.NewUi(asts, uiMessenger, controls, components, menuTree)
	onError(err, "failed to create gui")

	// Debug FTPS display
	if debugMode {
		gui.ToggleFtpsDisplay()
	}

	// Run UI, blocking
	err = ebiten.RunGame(gui)
	onError(err, "failed to run gui")
}

func onError(err error, msg string) {
	if err != nil {
		l := logger().With().Str("component", "main").Logger()
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
