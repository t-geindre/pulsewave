package main

import (
	"fmt"
	"synth/audio"
	"synth/instrument"
	intmidi "synth/midi"
	"time"

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
	pl.SetVolume(.4)
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

	for {
		time.Sleep(time.Second) // keep the program running
	}

	// Todo
	_ = stopFn
}
