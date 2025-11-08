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
	if len(os.Args) > 1 && os.Args[1] == "--debug" { // Todo implement flag parsing + buffer size option + device selection
		debugMode = true
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

	router.AddRoute(audioInQ, preset.ParamUpdateKind, uiOutQ)

	go router.Route()

	// Signal Chain
	synth := preset.NewPolysynth(SampleRate, audioOutQ, audioInQ)
	headroom := dsp.NewVca(synth, dsp.NewParam(0.8))
	clean := dsp.NewLowPassSVF(SampleRate, headroom, dsp.NewParam(18000), dsp.NewParam(0.5))

	// Midi setup
	midiListener := midi.NewListener(logger())
	defer midiListener.Close()

	device, err := midiListener.FindDevice()
	if err != nil {
		l := logger()
		l.Warn().Err(err).Msg("failed to find midiListener device")
		// TODO REMOVE ME
		go func() {
			for {
				midiInQ.TryWrite(msg.Message{
					Kind: midi.NoteOnKind,
					Key:  60,
				})
				time.Sleep(500 * time.Millisecond)
				midiInQ.TryWrite(msg.Message{
					Kind: midi.NoteOffKind,
					Key:  60,
				})
				time.Sleep(500 * time.Millisecond)
			}
		}()

	} else {
		err = midiListener.Listen(device, midiInQ)
		onError(err, "failed to listen to device")
	}

	// Player
	ctx := audio.NewContext(SampleRate)
	player, err := ctx.NewPlayerF32(dsp.NewStream(clean))
	onError(err, "failed to create player")

	player.SetBufferSize(time.Millisecond * 25)
	player.Play()

	// UI
	asts, err := assets.NewFromJson("assets/assets.json")
	onError(err, "failed to create assets loader")

	err = asts.Load()
	onError(err, "failed to load assets")

	ctrl := ui.NewMultiControls(
		ui.NewKeyboardControls(),
		ui.NewMidiControls(uiOutQ),
	)
	tree := preset.NewTree(uiOutQ, uiInQ)

	gui, err := ui.NewUi(asts, ctrl, tree)
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
