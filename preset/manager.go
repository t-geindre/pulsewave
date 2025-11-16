package preset

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"synth/dsp"
	"synth/msg"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

type presetVoice struct {
	preset *Preset
	voice  *Polysynth
	file   string
}

// todo because our goal is to be able to play multiple presets at once,
// we should have a way to unload presets without affecting others.
type Manager struct {
	*dsp.Mixer
	messenger *msg.Messenger
	current   int
	voices    []*presetVoice
	logger    zerolog.Logger
}

func NewManager(sr float64, logger zerolog.Logger, messenger *msg.Messenger, path string) (*Manager, error) {
	m := &Manager{
		Mixer:     dsp.NewMixer(dsp.NewParam(1), false),
		messenger: messenger,
		logger:    logger,
	}

	err := m.buildFromPath(sr, path)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) NoteOn(key int, vel float32) {
	m.voices[m.current].voice.NoteOn(key, vel)
}

func (m *Manager) NoteOff(key int) {
	m.voices[m.current].voice.NoteOff(key)
}

func (m *Manager) SetPitchBend(st float32) {
	m.voices[m.current].voice.SetPitchBend(st)
}

func (m *Manager) GetPresets() []string {
	names := make([]string, len(m.voices))
	for i, v := range m.voices {
		names[i] = v.preset.Name
	}
	return names
}

func (m *Manager) HandleMessage(msg msg.Message) {
	switch msg.Kind {
	case ParamUpdateKind:
		m.voices[m.current].voice.SetParam(msg.Key, msg.ValF)
	case LoadSavePresetKind:
		p := int(msg.Key)
		if msg.ValF == 0 {
			m.LoadPreset(p)
		} else if msg.ValF == 1 {
			m.SavePreset(p)
		}
	}
}

func (m *Manager) LoadPreset(p int) {
	if p < 0 || p >= len(m.voices) {
		return
	}

	// todo if not the current preset, silence notes from previous preset
	// todo send parameters updates to UI for the new preset

	m.voices[p].voice.LoadPreset(m.voices[p].preset) // reload preset
	m.current = p

	for key, param := range m.voices[p].preset.Params {
		m.messenger.SendMessage(msg.Message{
			Kind: ParamUpdateKind,
			Key:  key,
			ValF: param.GetBase(),
		})
	}
}

func (m *Manager) SavePreset(p int) {
	// Todo this may need to be done asynchronously to avoid blocking the audio thread

	if p < 0 || p >= len(m.voices) {
		return
	}

	preset := m.voices[p].preset
	m.voices[p].voice.HydratePreset(preset)

	prt := preset.ToProto()
	raw, err := proto.Marshal(prt)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to marshal preset")
		return
	}

	err = os.WriteFile(m.voices[p].file, raw, 0644)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to write preset file")
		return
	}

	m.logger.Info().
		Str("preset", preset.Name).
		Str("file", m.voices[p].file).
		Msg("preset saved")
}

func (m *Manager) buildFromPath(sr float64, pth string) error {
	files, err := filepath.Glob(filepath.Join(pth, "*.preset"))
	if err != nil {
		return err
	}

	voices := make([]*presetVoice, 0)
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		prt := &ProtoPreset{}
		err = proto.Unmarshal(raw, prt)
		if err != nil {
			return err
		}

		voice := &presetVoice{
			preset: NewPresetFromProto(prt),
			voice:  NewPolysynth(sr),
			file:   f,
		}

		voice.voice.LoadPreset(voice.preset)
		m.Mixer.Add(dsp.NewInput(voice.voice, nil, nil))

		voices = append(voices, voice)

		m.logger.Info().
			Str("preset", voice.preset.Name).
			Str("file", f).
			Msg("preset loaded")
	}

	slices.SortFunc(voices, func(a, b *presetVoice) int {
		return strings.Compare(a.preset.Name, b.preset.Name)
	})

	m.voices = voices

	// todo if no presets found, create a default one

	return nil
}
