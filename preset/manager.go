package preset

import (
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"synth/dsp"
	"synth/msg"
	"synth/settings"

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
	settings  map[uint8]dsp.Param
}

func NewManager(sr float64, logger zerolog.Logger, messenger *msg.Messenger, path string) *Manager {
	sets := make(map[uint8]dsp.Param)
	sets[settings.SettingsMasterGain] = dsp.NewSmoothedParam(sr, 1, 0.01)

	m := &Manager{
		Mixer:     dsp.NewMixer(sets[settings.SettingsMasterGain], false),
		messenger: messenger,
		logger:    logger,
		settings:  sets,
	}

	m.buildFromPath(sr, path)
	m.loadPreset(0) // force publish

	return m
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
	case PresetUpdateKind:
		m.voices[m.current].voice.SetParam(msg.Key, msg.ValF)
	case LoadSavePresetKind:
		p := int(msg.Key)
		if msg.ValF == 0 {
			m.loadPreset(p)
		} else if msg.ValF == 1 {
			m.savePreset(p)
		}
	case settings.SettingUpdateKind:
		if param, ok := m.settings[msg.Key]; ok {
			param.SetBase(msg.ValF)
		}
	}
}

func (m *Manager) loadPreset(p int) {
	if p < 0 || p >= len(m.voices) {
		return
	}

	if m.current != p {
		m.voices[m.current].voice.AllNotesOff()
	}

	m.voices[p].voice.LoadPreset(m.voices[p].preset) // reload preset
	m.current = p

	for key, param := range m.voices[p].preset.Params {
		m.messenger.SendMessage(msg.Message{
			Kind: PresetUpdateKind,
			Key:  key,
			ValF: param.GetBase(),
		})
	}
}

func (m *Manager) savePreset(p int) {
	// Todo this may need to be done asynchronously to avoid blocking the audio thread

	if p < 0 || p >= len(m.voices) {
		return
	}

	preset := m.voices[p].preset
	m.voices[p].voice.HydratePreset(preset)

	logger := m.logger.With().
		Str("preset", preset.Name).
		Str("file", m.voices[p].file).
		Logger()

	prt := preset.ToProto()
	raw, err := proto.Marshal(prt)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal preset")
		return
	}

	err = os.WriteFile(m.voices[p].file, raw, 0644)
	if err != nil {
		logger.Error().Err(err).Msg("failed to write preset file")
		return
	}

	logger.Info().Msg("preset saved")
}

func (m *Manager) buildFromPath(sr float64, pth string) {
	files, err := filepath.Glob(filepath.Join(pth, "*.preset"))
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to glob preset files")
		return
	}

	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			m.logger.Error().Err(err).Str("file", f).Msg("failed to read preset file")
			continue
		}

		prt := &ProtoPreset{}
		err = proto.Unmarshal(raw, prt)
		if err != nil {
			m.logger.Error().Err(err).Str("file", f).Msg("failed to unmarshal preset file")
			continue
		}

		preset := NewPresetFromProto(prt)
		m.addVoice(preset, sr, f)

		m.logger.Info().
			Str("preset", preset.Name).
			Str("file", f).
			Msg("preset loaded")
	}

	slices.SortFunc(m.voices, func(a, b *presetVoice) int {
		return strings.Compare(a.preset.Name, b.preset.Name)
	})

	if len(m.voices) > 0 {
		return
	}

	m.logger.Warn().Msg("no presets loaded, creating a default one")
	m.addVoice(NewPreset(), sr, path.Join(pth, "default.preset"))
}

func (m *Manager) addVoice(preset *Preset, sr float64, file string) {
	voice := &presetVoice{
		preset: preset,
		voice:  NewPolysynth(sr),
		file:   file,
	}
	voice.voice.LoadPreset(preset)
	m.Mixer.Add(dsp.NewInput(voice.voice, nil, nil))
	m.voices = append(m.voices, voice)
}
