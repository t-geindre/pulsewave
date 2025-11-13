package ui

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// touchClickMaxMove is the maximum distance in pixels for a press to be considered a click instead of a drag.
	touchClickMaxMove = 10
	// touchClickMaxDurationMs is the maximum duration in milliseconds for a press to be considered a click.
	touchClickMaxDurationMs = 250
	// touchDoubleClickMaxDelayMs is the maximum delay in milliseconds between two clicks to treat them as a double click.
	touchDoubleClickMaxDelayMs = 300
	// touchDoubleClickMaxDistance is the maximum distance in pixels between two clicks to still be a double click.
	touchDoubleClickMaxDistance = 20
	// touchSwipeMinDistance is the minimum horizontal distance in pixels to consider a gesture as a swipe.
	touchSwipeMinDistance = 60
	// touchSwipeMaxDurationMs is the maximum duration in milliseconds for a gesture to be considered a fast swipe.
	touchSwipeMaxDurationMs = 250
	// touchSwipeMaxOffAxisMove is the maximum vertical deviation in pixels for a gesture still considered horizontal.
	touchSwipeMaxOffAxisMove = 30
	// touchSwipeHorizontalDominanceFactor is the minimum ratio of horizontal over vertical movement to classify a swipe as horizontal.
	touchSwipeHorizontalDominanceFactor = 2.0
	// touchVelocitySampleCount is the number of recent movement samples used to estimate scroll inertia velocity.
	touchVelocitySampleCount = 5
	// touchInertiaFriction is the multiplicative factor applied each frame to reduce inertia velocity.
	touchInertiaFriction = 0.95
	// touchInertiaMinVelocityStart is the minimum absolute velocity required to start inertia after a drag.
	touchInertiaMinVelocityStart = 0.5
	// touchInertiaStopThreshold is the minimum absolute velocity below which inertia is stopped.
	touchInertiaStopThreshold = 0.2
	// touchInertiaCarryOverFactor is the proportion of previous inertia velocity reinjected into new inertia after a drag.
	touchInertiaCarryOverFactor = 0.5
	// touchInertiaBoostMinDrag is the minimum total vertical drag distance in pixels required to reuse previous inertia on a new drag.
	touchInertiaBoostMinDrag = 15
	// touchInertiaMaxIdleMs is the maximum time in milliseconds since the last movement during which inertia can still start on release.
	touchInertiaMaxIdleMs = 80
	// touchPrecisionMaxDrag is the maximum total vertical drag distance in pixels under which inertia is never applied to allow precise adjustments.
	touchPrecisionMaxDrag = 100
	// touchPrecisionStepPixels is the number of vertical pixels per one unit of scroll in the precision zone.
	touchPrecisionStepPixels = 30
	// touchPrecisionIdleResetMs is the idle time in milliseconds after which the precision zone is re-centered to allow fine movements again.
	touchPrecisionIdleResetMs = 80
	// touchPrecisionMaxStep is the maximum per-frame vertical delta in pixels considered as precision movement using discrete steps.
	touchPrecisionMaxStep = 8
	// touchFastMaxStep is the maximum per-frame vertical delta in pixels allowed in fast scroll mode to keep motion controllable.
	touchFastMaxStep = 24
	// touchFastSlope is the factor controlling how fast the fast scroll speed grows with additional pointer speed beyond the precision threshold.
	touchFastSlope = 0.15
	// touchSwipeScrollSuppressMinDx is the minimum horizontal drag distance in pixels after which vertical scrolling is suppressed for swipe gestures.
	touchSwipeScrollSuppressMinDx = 20
)

type TouchControls struct {
	isDown    bool
	startX    int
	startY    int
	lastX     int
	lastY     int
	startTime time.Time
	moved     bool

	pendingClick     bool
	pendingClickTime time.Time
	pendingClickX    int
	pendingClickY    int

	velSamples     [touchVelocitySampleCount]int
	velSampleCount int
	velSampleIndex int

	inertiaVelocity     float64
	prevInertiaVelocity float64

	totalDragDx int
	totalDragDy int

	precisionAcc           int
	lastMoveTime           time.Time
	suppressVerticalScroll bool
}

func NewTouchControls() *TouchControls {
	return &TouchControls{}
}

func (c *TouchControls) Update() (horDelta, vertDelta int) {
	horDelta, vertDelta = 0, 0
	now := time.Now()
	isDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	x, y := ebiten.CursorPosition()

	if !isDown && c.pendingClick && now.Sub(c.pendingClickTime).Milliseconds() > int64(touchDoubleClickMaxDelayMs) {
		horDelta = 1
		c.pendingClick = false
	}

	if isDown {
		if !c.isDown {
			c.isDown = true
			c.startX, c.startY = x, y
			c.lastX, c.lastY = x, y
			c.startTime = now
			c.lastMoveTime = now
			c.moved = false
			c.totalDragDx = 0
			c.totalDragDy = 0
			c.precisionAcc = 0
			c.suppressVerticalScroll = false

			c.prevInertiaVelocity = c.inertiaVelocity
			c.inertiaVelocity = 0

			c.velSampleCount = 0
			c.velSampleIndex = 0
		} else {
			dx := x - c.lastX
			dy := y - c.lastY

			if dx == 0 && dy == 0 {
				c.addVelocitySample(0)
			} else {
				if distanceSquared(c.startX, c.startY, x, y) > touchClickMaxMove*touchClickMaxMove {
					c.moved = true
				}

				c.totalDragDx += dx
				c.totalDragDy += dy

				if !c.suppressVerticalScroll {
					if absInt(c.totalDragDx) >= touchSwipeScrollSuppressMinDx &&
						float64(absInt(c.totalDragDx)) >= touchSwipeHorizontalDominanceFactor*float64(absInt(c.totalDragDy)) {
						c.suppressVerticalScroll = true
						c.precisionAcc = 0
					}
				}

				idleMs := now.Sub(c.lastMoveTime).Milliseconds()
				if idleMs > int64(touchPrecisionIdleResetMs) {
					c.precisionAcc = 0
				}

				vertDelta = 0

				if c.suppressVerticalScroll {
					c.addVelocitySample(0)
				} else {
					if absInt(dy) <= touchPrecisionMaxStep {
						c.precisionAcc += dy
						if c.precisionAcc >= touchPrecisionStepPixels {
							vertDelta = 1
							c.precisionAcc -= touchPrecisionStepPixels
						} else if c.precisionAcc <= -touchPrecisionStepPixels {
							vertDelta = -1
							c.precisionAcc += touchPrecisionStepPixels
						}
					} else {
						sign := 1
						if dy < 0 {
							sign = -1
						}
						mag := absInt(dy)
						excess := mag - touchPrecisionMaxStep
						speed := float64(touchPrecisionMaxStep) + float64(excess)*touchFastSlope
						if speed > float64(touchFastMaxStep) {
							speed = float64(touchFastMaxStep)
						}
						vertDelta = sign * int(math.Round(speed))
						c.precisionAcc = 0
					}

					if vertDelta != 0 {
						c.addVelocitySample(vertDelta)
					} else {
						c.addVelocitySample(0)
					}
				}

				c.lastX, c.lastY = x, y
				c.lastMoveTime = now
			}
		}
	} else {
		if c.isDown {
			c.isDown = false
			totalDx := c.totalDragDx
			totalDy := c.totalDragDy
			durMs := now.Sub(c.startTime).Milliseconds()
			dist2 := distanceSquared(c.startX, c.startY, x, y)
			idleMs := now.Sub(c.lastMoveTime).Milliseconds()

			if dist2 <= touchClickMaxMove*touchClickMaxMove && durMs <= int64(touchClickMaxDurationMs) {
				if c.pendingClick &&
					now.Sub(c.pendingClickTime).Milliseconds() <= int64(touchDoubleClickMaxDelayMs) &&
					distanceSquared(c.pendingClickX, c.pendingClickY, x, y) <= touchDoubleClickMaxDistance*touchDoubleClickMaxDistance {
					horDelta = -1
					c.pendingClick = false
				} else {
					c.pendingClick = true
					c.pendingClickTime = now
					c.pendingClickX, c.pendingClickY = x, y
				}
				c.prevInertiaVelocity = 0
				c.inertiaVelocity = 0
			} else {
				if durMs <= int64(touchSwipeMaxDurationMs) &&
					absInt(totalDx) >= touchSwipeMinDistance &&
					absInt(totalDy) <= touchSwipeMaxOffAxisMove &&
					float64(absInt(totalDx)) >= touchSwipeHorizontalDominanceFactor*float64(absInt(totalDy)) {
					if totalDx < 0 {
						horDelta = 1
					} else {
						horDelta = -1
					}
					c.prevInertiaVelocity = 0
					c.inertiaVelocity = 0
				} else if c.suppressVerticalScroll {
					c.prevInertiaVelocity = 0
					c.inertiaVelocity = 0
				} else {
					if absInt(c.totalDragDy) <= touchPrecisionMaxDrag {
						c.prevInertiaVelocity = 0
						c.inertiaVelocity = 0
					} else if idleMs > int64(touchInertiaMaxIdleMs) {
						c.prevInertiaVelocity = 0
						c.inertiaVelocity = 0
					} else {
						v := c.estimateVelocity()
						vTotal := v
						if absInt(c.totalDragDy) >= touchInertiaBoostMinDrag {
							vTotal += c.prevInertiaVelocity * touchInertiaCarryOverFactor
						}
						c.prevInertiaVelocity = 0
						if math.Abs(vTotal) >= touchInertiaMinVelocityStart {
							c.inertiaVelocity = vTotal
						} else {
							c.inertiaVelocity = 0
						}
					}
				}
			}

			c.velSampleCount = 0
			c.velSampleIndex = 0
			c.precisionAcc = 0
			c.suppressVerticalScroll = false
			c.totalDragDx = 0
			c.totalDragDy = 0
		}

		if !c.isDown && c.inertiaVelocity != 0 {
			vertDelta += int(math.Round(c.inertiaVelocity))
			c.inertiaVelocity *= touchInertiaFriction
			if math.Abs(c.inertiaVelocity) < touchInertiaStopThreshold {
				c.inertiaVelocity = 0
			}
		}
	}

	return horDelta, vertDelta
}

func (c *TouchControls) addVelocitySample(d int) {
	if touchVelocitySampleCount == 0 {
		return
	}
	c.velSamples[c.velSampleIndex] = d
	if c.velSampleCount < touchVelocitySampleCount {
		c.velSampleCount++
	}
	c.velSampleIndex = (c.velSampleIndex + 1) % touchVelocitySampleCount
}

func (c *TouchControls) estimateVelocity() float64 {
	if c.velSampleCount == 0 {
		return 0
	}
	sum := 0
	for i := 0; i < c.velSampleCount; i++ {
		sum += c.velSamples[i]
	}
	return float64(sum) / float64(c.velSampleCount)
}

func distanceSquared(x1, y1, x2, y2 int) int {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
