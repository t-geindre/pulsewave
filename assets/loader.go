package assets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	TypeImage = iota
	TypeFont
	TypeFace
	TypeRaw
)

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

func (l *Loader) MustLoad() {
	if l.Loaded {
		panic("loader already loaded")
	}

	for _, item := range l.ToLoad {
		switch item.Type {
		case TypeImage:
			l.Images[item.Name] = MustLoadImage(item.Path)
		case TypeFont:
			l.Fonts[item.Name] = MustLoadFont(item.Path)
		case TypeRaw:
			l.Raws[item.Name] = MustLoadRaw(item.Path)
		case TypeFace: //Depends on font being loaded first
		default:
			panic("Unknown asset type")
		}
	}

	for _, item := range l.ToLoad {
		if item.Type == TypeFace {
			font, ok := l.Fonts[item.Path]
			if !ok {
				panic("Font not found for face: " + item.Name)
			}
			l.Faces[item.Name] = &text.GoTextFace{
				Source: font,
				Size:   item.Val,
			}
		}
	}

	l.ToLoad = nil
	l.Loaded = true
}

func (l *Loader) GetImage(name string) *ebiten.Image {
	img, ok := l.Images[name]
	if !ok {
		panic("Image not found: " + name)
	}

	return img
}

func (l *Loader) GetFont(name string) *text.GoTextFaceSource {
	font, ok := l.Fonts[name]
	if !ok {
		panic("Font not found: " + name)
	}

	return font
}

func (l *Loader) GetFace(name string) text.Face {
	face, ok := l.Faces[name]
	if !ok {
		panic("Face not found: " + name)
	}

	return face
}

func (l *Loader) GetRaw(name string) []byte {
	raw, ok := l.Raws[name]
	if !ok {
		panic("Raw data not found: " + name)
	}

	return raw
}

func (l *Loader) GetPath(name string) string {
	for _, item := range l.ToLoad {
		if item.Name == name {
			return item.Path
		}
	}
	panic("Path not found for asset: " + name)
}
