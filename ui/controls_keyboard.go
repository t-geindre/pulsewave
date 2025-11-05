package ui

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyboardControls struct {
	// Vertical keys with acceleration
	vertKeys      map[ebiten.Key]int
	vertKeysAccel map[time.Duration]int
	vertKeysStart map[ebiten.Key]time.Time

	// Horizontal keys (no acceleration)
	horKeys map[ebiten.Key]int
}

func NewKeyboardControls() *KeyboardControls {
	return &KeyboardControls{
		vertKeysAccel: map[time.Duration]int{
			200 * time.Millisecond:  1,
			500 * time.Millisecond:  1,
			1000 * time.Millisecond: 5,
			1500 * time.Millisecond: 5,
			2000 * time.Millisecond: 10,
			2500 * time.Millisecond: 20,
			3500 * time.Millisecond: 40,
		},
		vertKeys: map[ebiten.Key]int{
			ebiten.KeyUp:   -1,
			ebiten.KeyDown: 1,
		},
		vertKeysStart: make(map[ebiten.Key]time.Time),
		horKeys: map[ebiten.Key]int{
			ebiten.KeyLeft:  -1,
			ebiten.KeyRight: 1,
		},
	}
}

func (c *KeyboardControls) Update() (int, int) {
	horDelta, vertDelta := 0, 0

	for k, dir := range c.horKeys {
		if inpututil.IsKeyJustPressed(k) {
			horDelta += dir
		}
	}

	for k, dir := range c.vertKeys {
		if ebiten.IsKeyPressed(k) {
			if _, ok := c.vertKeysStart[k]; !ok {
				c.vertKeysStart[k] = time.Now()
				vertDelta += dir
			} else {
				elapsed := time.Since(c.vertKeysStart[k])
				for d, v := range c.vertKeysAccel {
					if elapsed > d {
						vertDelta += v * dir
					}
				}
			}
		} else {
			delete(c.vertKeysStart, k)
		}
	}

	return horDelta, vertDelta
}
