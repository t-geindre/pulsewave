package assets

import (
	"errors"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	TypeImage = iota
	TypeFont
	TypeFace
	TypeRaw
)

var ErrUnknownAssetType = errors.New("unknown asset type")
var ErrFontNotFound = errors.New("font not found for face")
var ErrAssetNotFound = errors.New("asset not found")

type toLoad struct {
	Path, Name string
	Val        float64
	Type       int
}

type Loader struct {
	Images map[string]*ebiten.Image
	Fonts  map[string]*text.GoTextFaceSource
	Raws   map[string][]byte
	Faces  map[string]text.Face

	ToLoad []toLoad
	Loaded bool
}

func NewLoader() *Loader {
	return &Loader{
		Images: make(map[string]*ebiten.Image),
		Fonts:  make(map[string]*text.GoTextFaceSource),
		Raws:   make(map[string][]byte),
		Faces:  make(map[string]text.Face),
	}
}

func (l *Loader) AddImage(name, path string) {
	l.Add(name, path, TypeImage, 0)
}

func (l *Loader) AddFont(name, path string) {
	l.Add(name, path, TypeFont, 0)
}

func (l *Loader) AddFace(name, font string, size float64) {
	l.Add(name, font, TypeFace, size)
}

func (l *Loader) AddRaw(name, path string) {
	l.Add(name, path, TypeRaw, 0)
}

func (l *Loader) Add(name, path string, t int, v float64) {
	l.ToLoad = append(l.ToLoad, toLoad{
		Name: name,
		Path: path,
		Type: t,
		Val:  v,
	})
}

func (l *Loader) Load() error {
	if l.Loaded {
		return nil
	}

	for _, item := range l.ToLoad {
		switch item.Type {
		case TypeImage:
			img, _, err := ebitenutil.NewImageFromFile(item.Path)
			if err != nil {
				return err
			}
			l.Images[item.Name] = img
		case TypeFont:
			obj, err := doLoadFont(item.Path)
			if err != nil {
				return err
			}
			l.Fonts[item.Name] = obj
		case TypeRaw:
			raw, err := os.ReadFile(item.Path)
			if err != nil {
				return err
			}
			l.Raws[item.Name] = raw
		case TypeFace: //Depends on font being loaded first
		default:
			return ErrUnknownAssetType
		}
	}

	for _, item := range l.ToLoad {
		if item.Type == TypeFace {
			font, ok := l.Fonts[item.Path]
			if !ok {
				return fmt.Errorf("%s: %w", item.Path, ErrFontNotFound)
			}
			l.Faces[item.Name] = &text.GoTextFace{
				Source: font,
				Size:   item.Val,
			}
		}
	}

	l.ToLoad = nil
	l.Loaded = true

	return nil
}

func (l *Loader) GetImage(name string) (*ebiten.Image, error) {
	img, ok := l.Images[name]
	if !ok {
		return nil, fmt.Errorf("%v: Image %s", ErrAssetNotFound, name)
	}

	return img, nil
}

func (l *Loader) GetFont(name string) (*text.GoTextFaceSource, error) {
	font, ok := l.Fonts[name]
	if !ok {
		return nil, fmt.Errorf("%v: Font %s", ErrAssetNotFound, name)
	}

	return font, nil
}

func (l *Loader) GetFace(name string) (text.Face, error) {
	face, ok := l.Faces[name]
	if !ok {
		return nil, fmt.Errorf("%v: Face %s", ErrAssetNotFound, name)
	}

	return face, nil
}

func (l *Loader) GetRaw(name string) ([]byte, error) {
	raw, ok := l.Raws[name]
	if !ok {
		return nil, fmt.Errorf("%v:tRaw %s", ErrAssetNotFound, name)
	}

	return raw, nil
}

func doLoadFont(path string) (*text.GoTextFaceSource, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	f, err := text.NewGoTextFaceSource(r)
	if err != nil {
		return nil, err
	}

	return f, nil
}
