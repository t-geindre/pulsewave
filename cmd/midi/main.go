package main

import (
	"fmt"
	"synth/assets"
	"synth/audio"
	"synth/instrument"
	intmidi "synth/midi"
	"synth/ui"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	// SOUND SETUP
	const SampleRate = 44100

	voicer := intmidi.NewVoicer(8, func() audio.Source {
		return instrument.NewSynth(SampleRate)
	})

	pl := audio.NewPlayer(SampleRate, voicer)
	pl.SetVolume(.8)
	pl.SetBufferSize(time.Millisecond * 40)

	// ROUTING SETUP
	router := intmidi.NewRouter(voicer, intmidi.NewMenu())

	// MIDI SETUP
	defer midi.CloseDriver()

	ips := midi.GetInPorts()
	in, err := midi.FindInPort(ips[1].String())
	if err != nil {
		fmt.Println("can't find VMPK")
		return
	}

	stopFn, err := midi.ListenTo(in, func(msg midi.Message, _ int32) {
		router.Route(msg)
	})

	if err != nil {
		panic(err)
	}

	// Todo
	_ = stopFn

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
