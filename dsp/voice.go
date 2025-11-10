package dsp

type Resettable interface {
	Reset(bool)
}

type Envelope interface {
	NoteOn()
	NoteOff()
	IsIdle() bool
	ParamModulator
}

type Voice struct {
	Node
	freq   Param
	envs   []Envelope
	resets []Resettable
}

func NewVoice(src Node, freq Param, extra ...any) *Voice {
	v := &Voice{
		Node:   src,
		freq:   freq,
		envs:   make([]Envelope, 0),
		resets: make([]Resettable, 0),
	}

	for _, e := range extra {
		switch e := e.(type) {
		case Envelope:
			v.envs = append(v.envs, e)
		case Resettable:
			v.resets = append(v.resets, e)
		default:
			panic("invalid extra")
		}
	}

	return v
}

func (v *Voice) NoteOn(key int, vel float32) {
	// v.gain.SetBase(vel) todo handle vel, probably with a param modulator
	v.freq.SetBase(MidiKeys[key])
	soft := !v.envs[0].IsIdle()
	v.Node.Reset(soft)

	for _, reset := range v.resets {
		reset.Reset(soft)
	}

	for _, env := range v.envs {
		env.NoteOn()
	}
}

func (v *Voice) NoteOff() {
	for _, env := range v.envs {
		env.NoteOff()
	}
}

func (v *Voice) IsIdle() bool {
	for _, env := range v.envs {
		if !env.IsIdle() {
			return false
		}
	}
	return true
}
