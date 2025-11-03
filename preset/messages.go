package preset

import "synth/msg"

const AudioSource msg.Source = 10

const ParamUpdateKind msg.Kind = 1

const FBDelayParam = iota

/*
Message
	Source Source AudioSource
	Kind   Kind   ParamUpdateKind
	Key    uint8  0-255 parameter ID
	Val16  int16  0-65535
}
*/
