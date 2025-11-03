package midi

import "synth/msg"

const MidiSource msg.Source = 0

const NoteOnKind msg.Kind = 1
const NoteOffKind msg.Kind = 2
const PitchBendKind msg.Kind = 4
const ControlChangeKind msg.Kind = 3
