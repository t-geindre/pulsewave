package audio

import "github.com/hajimehoshi/ebiten/v2/audio"

type Player struct {
	audio.Player
	src Source
}

func NewPlayer(sr int, src Source) *Player {
	ado := audio.NewContext(sr)
	player, err := ado.NewPlayerF32(NewStream(src))
	if err != nil {
		panic(err)
	}

	player.SetVolume(1)
	player.Play()

	return &Player{
		Player: *player,
		src:    src,
	}
}

func (p *Player) IsPlaying() bool {
	return p.src.IsActive()
}
