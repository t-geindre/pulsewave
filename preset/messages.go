package preset

import "synth/msg"

const AudioSource msg.Source = 10

const ParamUpdateKind msg.Kind = 1
const ParamPullAllKind msg.Kind = 2

const FBDelayParam = iota

/*
Message
	Source Source AudioSource
	Kind   Kind   ParamUpdateKind
	key    uint8  0-255 parameter ID
	ValF  float32
}
*/
