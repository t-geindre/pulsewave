package preset

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

var ErrRootIsNotTree = errors.New("cannot load preset, root is not *Tree")

type PresetNode struct {
	preset *Preset
	val    float32
	file   string
	logger zerolog.Logger
	*SelectorNode
}

func NewPresetNode(file string, logger zerolog.Logger) *PresetNode {
	logger = logger.With().Str("component", "PresetNode").Str("file", file).Logger()

	data, err := os.ReadFile(file)
	if err != nil {
		logger.Warn().Err(err).Msg("cannot load preset")
	}
	pp := &ProtoPreset{}
	err = proto.Unmarshal(data, pp)
	if err != nil {
		logger.Err(err).Msg("cannot load preset")
	}

	logger = logger.With().Str("preset", pp.Name).Logger()

	preset := NewPresetFromProto(pp)
	if preset.Name == "" {
		preset.Name = "Unknown"
	}

	return &PresetNode{
		SelectorNode: NewSelectorNode(preset.Name, NONE,
			NewSelectorOption("Load", "", 0),
			NewSelectorOption("Save", "", 1),
		),
		preset: preset,
		file:   file,
		logger: logger,
	}
}

func (pn *PresetNode) Val() float32 {
	return pn.val
}

func (pn *PresetNode) SetVal(v float32) {
	pn.val = v
}

func (pn *PresetNode) Validate() {
	switch int(pn.val) {
	case 0:
		pn.Load()
	case 1:
		pn.Save()
	}
}

func (pn *PresetNode) Load() {
	tree, ok := pn.Root().(*Tree)
	if !ok {
		pn.logger.Fatal().Err(ErrRootIsNotTree).Msg("cannot load preset")
	}

	tree.LoadPreset(pn.preset)
	pn.logger.Info().Msg("preset loaded")
}

func (pn *PresetNode) Save() {
	tree, ok := pn.Root().(*Tree)
	if !ok {
		pn.logger.Fatal().Err(ErrRootIsNotTree).Msg("cannot save preset")
	}

	preset := tree.GetPreset()
	preset.Name = pn.preset.Name
	pn.preset = preset

	data, err := proto.Marshal(preset.ToProto())
	if err != nil {
		pn.logger.Err(err).Msg("cannot save preset")
	}

	err = os.WriteFile(pn.file, data, 0644)
	if err != nil {
		pn.logger.Err(err).Msg("cannot save preset")
	}

	pn.logger.Info().Msg("preset saved")
}
