package main

import (
	"synth/audio"
	"synth/song"
	"time"
)

const SampleRate = 44100

func main() {
	player := audio.NewPlayer(SampleRate, song.NewCrazyFrog(SampleRate))

	for player.IsPlaying() {
		time.Sleep(100 * time.Millisecond)
	}
}
