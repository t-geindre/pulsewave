package assets

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func MustLoadRaw(path string) []byte {
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return raw
}

func MustLoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)

	if err != nil {
		panic(err)
	}

	return img
}

func MustLoadShader(path string) *ebiten.Shader {
	shader, err := ebiten.NewShader(MustLoadRaw(path))
	if err != nil {
		panic(err)
	}

	return shader
}

func MustLoadFont(path string) *text.GoTextFaceSource {
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	f, err := text.NewGoTextFaceSource(r)
	if err != nil {
		panic(err)
	}

	return f
}
