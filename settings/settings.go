package settings

import (
	"os"
	"sync"
	"synth/msg"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

type Settings struct {
	settings  map[uint8]float32
	path      string
	logger    zerolog.Logger
	messenger *msg.Messenger
	closed    chan struct{}
	lock      sync.Mutex
	dirty     bool
}

// NewSettings creates a new Settings manager.
// Loads settings from the given file path.
func NewSettings(pth string, messenger *msg.Messenger, logger zerolog.Logger) *Settings {
	s := &Settings{
		settings:  make(map[uint8]float32),
		path:      pth,
		logger:    logger.With().Str("file", pth).Logger(),
		messenger: messenger,
		closed:    make(chan struct{}),
	}

	s.loadDefaults()
	s.load()

	messenger.RegisterHandler(s)

	go s.periodicPersist()

	return s
}

// Set sets a setting value by its ID.
// Thread-safe.
func (s *Settings) Set(id uint8, value float32) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if setting, ok := s.settings[id]; ok {
		if setting != value {
			s.settings[id] = value
			s.dirty = true
		}
	}

	s.messenger.SendMessage(msg.Message{
		Kind: SettingUpdateKind,
		Key:  id,
		ValF: value,
	})
}

// HandleMessage processes incoming messages to update settings.
func (s *Settings) HandleMessage(msg msg.Message) {
	if msg.Kind != SettingUpdateKind {
		return
	}

	s.Set(msg.Key, msg.ValF)
}

// Persist saves the current settings to the file.
// Thread-safe.
func (s *Settings) Persist() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.dirty {
		return
	}
	s.dirty = false

	prt := &ProtoSettings{}
	for id, value := range s.settings {
		prt.Settings = append(prt.Settings, &ProtoSetting{
			Id:    uint32(id),
			Value: value,
		})
	}

	raw, err := proto.Marshal(prt)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal settings")
		return
	}

	err = os.WriteFile(s.path, raw, 0644)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to write settings file")
		return
	}

	s.logger.Info().Msg("settings persisted")
}

// Close stops the periodic persistence goroutine.
func (s *Settings) Close() {
	close(s.closed)
}

func (s *Settings) load() {
	raw, err := os.ReadFile(s.path)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to read settings file, using default settings")
		return
	}

	prt := &ProtoSettings{}
	err = proto.Unmarshal(raw, prt)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to unmarshal settings file, using default settings")
		return
	}

	for id, value := range s.settings {
		found := false
		for _, setting := range prt.Settings {
			if uint8(setting.Id) == id {
				s.Set(uint8(setting.Id), setting.Value)
				found = true
				break
			}
		}
		if !found {
			s.Set(id, value)
		}
	}

	s.dirty = false // avoid useless persist right after load
	s.logger.Info().Msg("settings loaded")
}

func (s *Settings) loadDefaults() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.settings[MasterGain] = 1.0
	s.settings[PitchBendRange] = 4.0
}

func (s *Settings) periodicPersist() {
	persistTicker := time.NewTicker(time.Second * 10)
	defer persistTicker.Stop()

	updateTicker := time.NewTicker(time.Millisecond)
	defer updateTicker.Stop()

	for {
		select {
		case <-persistTicker.C:
			s.Persist()
		case <-s.closed:
			return
		case <-updateTicker.C:
			s.messenger.Process()
		}
	}
}
