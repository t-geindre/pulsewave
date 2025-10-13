package automation

import "time"

type NextFunc func(t time.Duration) (NextFunc, time.Duration)

type Timed struct {
	time time.Duration
	step time.Duration

	nextFunc NextFunc
	nextWhen time.Duration
}

func NewTimed(sr float64, first NextFunc) *Timed {
	return &Timed{
		time:     0,
		step:     time.Second / time.Duration(sr),
		nextFunc: first,
		nextWhen: 0,
	}
}

func (t *Timed) NextSample() float64 {
	if t.time >= t.nextWhen {
		t.nextFunc, t.nextWhen = t.nextFunc(t.time)
		t.time = 0
		return 0
	}

	t.time += t.step
	return 0
}
