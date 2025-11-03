package msg

type Kind uint8
type Source uint8

type Message struct {
	Kind   Kind
	Source Source
	Key    uint8 // note number, controller number, ...
	Val8   uint8 // velocity, control value, ...
	Val16  int16
	ValF   float32
	Chan   uint8
}
