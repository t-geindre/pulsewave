package msg

type Kind uint8
type Source uint8

type Message struct {
	Kind   Kind
	Source Source
	Key    uint8 // note number, controller number, ...
	Val    uint8 // velocity, control value, ...
	Chan   uint8
}
