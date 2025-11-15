package msg

type Kind uint8

type Message struct {
	Kind  Kind
	Key   uint8 // midi key, controller number, ...
	Chan  uint8 // midi channel
	Val8  uint8 // velocity, control value, ...
	Val16 int16 // pitch bend, ...
	ValF  float32
}
