package audio

import "github.com/hajimehoshi/ebiten/v2/audio"

type Player struct {
	audio.Player
}

func NewPlayer(sr int, src Source) *Player {
	ado := audio.NewContext(sr)
	player, err := ado.NewPlayerF32(NewStream(src))
	if err != nil {
		panic(err)
	}

	player.SetVolume(.5)
	player.Play()

	return &Player{*player}
}
