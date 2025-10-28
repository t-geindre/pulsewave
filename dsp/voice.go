package dsp

type Resettable interface {
	Reset()
}

type Envelope interface {
	NoteOn()
	NoteOff()
	IsIdle() bool
	ParamModulator
}

type Voice struct {
	Node
	freq  Param
	gain  Param
	env   Envelope
	extra []Resettable
}

func NewVoice(src Node, freq Param, env Envelope, extra ...Resettable) *Voice {
	gain := NewParam(0)
	*gain.ModInputs() = append(*gain.ModInputs(), NewModInput(env, 1.0, nil))

	return &Voice{
		Node:  NewVca(src, gain),
		freq:  freq,
		gain:  gain,
		env:   env,
		extra: extra,
	}
}

func (v *Voice) NoteOn(key int, vel float32) {
	// v.gain.SetBase(vel) todo handle vel, probably with a param modulator
	v.freq.SetBase(MidiKeys[key])

	if v.env.IsIdle() {
		v.Node.Reset()
	}
	for _, n := range v.extra {
		n.Reset()
	}

	v.env.NoteOn()
}

func (v *Voice) NoteOff() {
	v.env.NoteOff()
}

func (v *Voice) IsIdle() bool {
	return v.env.IsIdle()
}
