package main

import (
	"fmt"
	"math"
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

	voiceFact := func() *dsp.Voice {
		// Oscillators
		mixer := dsp.NewMixer(dsp.NewParam(1), true)
		reg := dsp.NewShapeRegistry()
		freq := dsp.NewSmoothedParam(SampleRate, 440, .001)
		oscg := float32(1.0 / math.Sqrt(3))

		// 0
		reg.Set(0, dsp.ShapeSaw)
		f0 := dsp.NewTunerParam(freq, dsp.NewParam(-0.07))
		ps0 := dsp.NewParam(.33)
		mixer.Add(dsp.NewInput(
			dsp.NewRegOscillator(SampleRate, reg, 0, f0, ps0, nil),
			dsp.NewParam(oscg),
			dsp.NewParam(-0.3),
		))

		// 1
		reg.Set(0, dsp.ShapeSaw)
		mixer.Add(dsp.NewInput(
			dsp.NewRegOscillator(SampleRate, reg, 0, freq, nil, nil),
			dsp.NewParam(oscg),
			dsp.NewParam(0),
		))

		// 2
		reg.Set(0, dsp.ShapeSaw)
		f1 := dsp.NewTunerParam(freq, dsp.NewParam(+0.07))
		mixer.Add(dsp.NewInput(
			dsp.NewRegOscillator(SampleRate, reg, 0, f1, dsp.NewParam(.66), nil),
			dsp.NewParam(oscg),
			dsp.NewParam(.3),
		))

		// LPF
		cutoff := dsp.NewSmoothedParam(SampleRate, 800, 0.005)
		reson := dsp.NewParam(1)
		lpf := dsp.NewLowPassSVF(SampleRate, mixer, cutoff, reson)

		ctModRateAdsr := dsp.NewADSR(SampleRate, 0, time.Millisecond*50, 0, time.Millisecond*100)
		*cutoff.ModInputs() = append(*cutoff.ModInputs(), dsp.NewModInput(ctModRateAdsr, 1000, nil))

		ctModRateOsc := dsp.NewOscillator(SampleRate, dsp.ShapeSine, dsp.NewParam(.5), dsp.NewParam(1), nil)
		*cutoff.ModInputs() = append(*cutoff.ModInputs(), dsp.NewModInput(ctModRateOsc, 200, nil))

		// Voice
		adsr := dsp.NewADSR(SampleRate, time.Millisecond*10, time.Millisecond*800, .9, time.Millisecond*100)
		voice := dsp.NewVoice(lpf, freq, adsr, ctModRateAdsr, ctModRateOsc)

		return voice
	}

	// Polyphonic voice
	poly := dsp.NewPolyVoice(8, voiceFact)

	// Delay
	delay := dsp.NewFeedbackDelay(
		SampleRate,
		2.0,
		poly,
		dsp.NewParam(0.35), // delay time in seconds
		dsp.NewParam(0.3),  // feedback amount (0-1)
		dsp.NewParam(0.2),  // mix
		dsp.NewParam(2000), // mix amount (0-1)
	)

	// Player
	p := audio.NewPlayer(SampleRate, delay)
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
			poly.NoteOn(int(key), float32(vel/127.0))
		case msg.GetNoteEnd(&ch, &key):
			poly.NoteOff(int(key))
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
