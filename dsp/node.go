package dsp

type Node interface {
	Source
	Reset(soft bool)
}
