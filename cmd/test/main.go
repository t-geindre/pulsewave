package main

import (
	"fmt"
	"synth/assets"
	"synth/audio"
	"synth/dsp"
	"synth/ui"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	const SampleRate = 44100

	// Mixer
	mixer := dsp.NewMixer(dsp.NewParam(1), true)

	// Oscillators
	reg := dsp.NewShapeRegistry()
	reg.Set(0, dsp.ShapeSine)
	reg.Set(1, dsp.ShapeSaw)
	reg.Set(2, dsp.ShapeTriangle)

	freq := dsp.NewSmoothedParam(SampleRate, 440, .001)
	for i := 0; i < 3; i++ {
		oscillator := dsp.NewRegOscillator(SampleRate, reg, i, freq, nil, nil)
		input := &dsp.Input{
			Src:  oscillator,
			Gain: dsp.NewParam(0.3),
			Pan:  dsp.NewParam(-0.5 + float32(i)*0.5),
		}
		mixer.Add(input)
	}

	// Freq mod
	//fmRate := dsp.NewOscillator(SampleRate, dsp.ShapeSine, dsp.NewParam(1), nil, nil)
	//*freq.ModInputs() = append(*freq.ModInputs(), dsp.NewModInput(fmRate, 100.0, nil))

	// Envelope
	gain := dsp.NewParam(0)
	adsr := dsp.NewADSR(SampleRate, time.Millisecond*10, time.Millisecond*100, 0.8, time.Millisecond*300)
	*gain.ModInputs() = append(*gain.ModInputs(), dsp.NewModInput(adsr, 1.0, nil))
	vca := dsp.NewVca(mixer, gain)

	// Player
	p := audio.NewPlayer(SampleRate, vca)
	p.SetBufferSize(time.Millisecond * 20)

	// MIDI SETUP
	defer midi.CloseDriver()

	ips := midi.GetInPorts()
	in, err := midi.FindInPort(ips[1].String())
	if err != nil {
		fmt.Println("can't find VMPK")
		return
	}

	stopFn, err := midi.ListenTo(in, func(msg midi.Message, _ int32) {
		var ch, key, vel uint8
		switch {
		case msg.GetNoteStart(&ch, &key, &vel):
			freq.SetBase(dsp.MidiKeys[key])
			adsr.NoteOn()
		case msg.GetNoteEnd(&ch, &key):
			adsr.NoteOff()
		default:
			fmt.Println(msg)
		}
	})
	_ = stopFn // Todo

	if err != nil {
		panic(err)
	}

	// UI
	asts := assets.NewLoader()
	asts.AddImage("ui/background", "assets/imgs/background.png")
	asts.AddImage("ui/arrow", "assets/imgs/arrow.png")
	asts.AddImage("ui/selected", "assets/imgs/selected.png")
	asts.AddFont("ui/font", "assets/fonts/roboto/Roboto-Medium.ttf") // Roboto medium, letter spacing 3, size 20, color white
	asts.AddFace("ui/face", "ui/font", 21)
	asts.MustLoad()

	err = ebiten.RunGame(ui.NewUi(asts))
	if err != nil {
		panic(err)
	}
}
