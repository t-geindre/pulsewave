package ui

import (
	"image/color"
	"synth/dsp"
	"synth/tree"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Oscilloscope struct {
	in *AudioQueue

	// Mono sample ring buffer
	ring    []float32
	mask    int
	write   int
	tmpMono []float32

	// Display params
	zoom float64 // samples per pixel
	gain float64 // amplitude scale

	// Trigger params
	trigger       bool
	triggerLevel  float32
	triggerRising bool
}

// NewOscilloscope expects a power of two ringSize
func NewOscilloscope(in *AudioQueue, ringSize int) (*Oscilloscope, error) {
	n := 1
	for n < ringSize {
		n <<= 1
	}
	return &Oscilloscope{
		in:            in,
		ring:          make([]float32, n),
		mask:          n - 1,
		tmpMono:       make([]float32, dsp.BlockSize),
		zoom:          1.0,
		gain:          1.8,
		trigger:       true,
		triggerRising: true,
	}, nil
}

// Update drains the audio queue into the ring buffer
func (o *Oscilloscope) Update() {
	if o.in == nil {
		return
	}

	o.in.Drain(0, func(b dsp.Block) {
		// Mix L/R to mono
		for i := 0; i < dsp.BlockSize; i++ {
			o.tmpMono[i] = 0.5 * (b.L[i] + b.R[i])
		}

		// Write samples to ring
		for i := 0; i < dsp.BlockSize; i++ {
			o.ring[o.write&o.mask] = o.tmpMono[i]
			o.write++
		}
	})
}

// Draw renders the oscilloscope waveform
func (o *Oscilloscope) Draw(screen *ebiten.Image) {
	if o.ring == nil || len(o.ring) == 0 {
		return
	}

	bds := screen.Bounds()
	w, h := bds.Dx(), bds.Dy()
	if w <= 1 || h <= 1 {
		return
	}

	// Integer samples-per-pixel step
	step := int(o.zoom + 0.5)
	if step < 1 {
		step = 1
	}

	// We draw w pixels => w-1 segments => (w-1)*step samples
	segments := w - 1
	totalSamples := step * segments
	if totalSamples <= 0 {
		return
	}

	// Clamp to ring size
	if totalSamples > len(o.ring) {
		segments = len(o.ring) / step
		if segments <= 0 {
			return
		}
		totalSamples = step * segments
	}

	mid := float64(h) / 2
	scale := mid * o.gain

	searchEnd := o.write // exclusive
	level := o.triggerLevel

	// Limit backwards search to ring size
	maxLookback := totalSamples * 2
	if maxLookback > len(o.ring) {
		maxLookback = len(o.ring)
	}
	searchStart := searchEnd - maxLookback

	// Last sample used must be <= searchEnd-1
	// windowStart + segments*step <= searchEnd-1
	lastPossibleStart := searchEnd - 1 - segments*step

	var windowStart int

	if o.trigger {
		found := false

		// Trigger must leave enough samples after it
		triggerSearchEnd := lastPossibleStart
		if triggerSearchEnd <= searchStart {
			triggerSearchEnd = searchStart
		}

		prev := o.sampleAt(searchStart)
		for pos := searchStart + 1; pos <= triggerSearchEnd; pos++ {
			cur := o.sampleAt(pos)

			if o.triggerRising {
				if prev < level && cur >= level {
					windowStart = pos
					found = true
					break
				}
			} else {
				if prev > level && cur <= level {
					windowStart = pos
					found = true
					break
				}
			}
			prev = cur
		}

		if !found {
			// Fallback to last window before end
			windowStart = lastPossibleStart
		}
	} else {
		windowStart = lastPossibleStart
	}

	// Draw waveform slightly outside the visible area
	for x := 0; x < segments; x++ {
		samplePos0 := windowStart + x*step
		samplePos1 := windowStart + (x+1)*step

		s0 := o.sampleAt(samplePos0)
		s1 := o.sampleAt(samplePos1)

		y0 := mid - float64(s0)*scale
		y1 := mid - float64(s1)*scale

		vector.StrokeLine(
			screen,
			float32(x), float32(y0),
			float32(x+1), float32(y1),
			3.0,
			color.White,
			true,
		)
	}
}

// sampleAt reads a logical sample from the ring
func (o *Oscilloscope) sampleAt(pos int) float32 {
	idx := pos & o.mask
	return o.ring[idx]
}

// Scroll updates the time zoom factor
func (o *Oscilloscope) Scroll(delta int) {
	if delta == 0 {
		return
	}

	step := 0.1

	if delta > 0 {
		o.zoom -= step * float64(delta)
	} else {
		o.zoom -= step * float64(delta)
	}

	if o.zoom < 0.1 {
		o.zoom = 0.1
	}
	if o.zoom > 10.0 {
		o.zoom = 10.0
	}
}

func (o *Oscilloscope) CurrentTarget() tree.Node {
	return nil
}

func (o *Oscilloscope) Focus() {}
func (o *Oscilloscope) Blur()  {}
