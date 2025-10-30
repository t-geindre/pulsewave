package msg

type Kind uint8
type Source uint8

type Message struct {
	Type       Kind
	Source     Source
	V1, V2, V3 uint8
}
